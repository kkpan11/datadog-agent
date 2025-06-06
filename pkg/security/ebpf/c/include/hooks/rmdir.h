#ifndef _HOOKS_RMDIR_H_
#define _HOOKS_RMDIR_H_

#include "constants/syscall_macro.h"
#include "helpers/approvers.h"
#include "helpers/events_predicates.h"
#include "helpers/filesystem.h"
#include "helpers/syscalls.h"

int __attribute__((always_inline)) trace__sys_rmdir(u8 async, const char *filename) {
    struct syscall_cache_t syscall = {
        .type = EVENT_RMDIR,
        .policy = fetch_policy(EVENT_RMDIR),
        .async = async,
    };

    if (!async) {
        collect_syscall_ctx(&syscall, SYSCALL_CTX_ARG_STR(0), (void *)filename, NULL, NULL);
    }
    cache_syscall(&syscall);

    return 0;
}

HOOK_SYSCALL_ENTRY1(rmdir, const char *, filename) {
    return trace__sys_rmdir(SYNC_SYSCALL, filename);
}

HOOK_ENTRY("do_rmdir")
int hook_do_rmdir(ctx_t *ctx) {
    struct syscall_cache_t *syscall = peek_syscall_with(rmdir_predicate);
    if (!syscall) {
        return trace__sys_rmdir(ASYNC_SYSCALL, NULL);
    }
    return 0;
}

// security_inode_rmdir is shared between rmdir and unlink syscalls
HOOK_ENTRY("security_inode_rmdir")
int hook_security_inode_rmdir(ctx_t *ctx) {
    struct syscall_cache_t *syscall = peek_syscall_with(rmdir_predicate);
    if (!syscall) {
        return 0;
    }

    struct path_key_t key = {};
    struct dentry *dentry = NULL;

    switch (syscall->type) {
    case EVENT_RMDIR:
        if (syscall->rmdir.file.path_key.ino) {
            return 0;
        }

        // we resolve all the information before the file is actually removed
        dentry = (struct dentry *)CTX_PARM2(ctx);
        set_file_inode(dentry, &syscall->rmdir.file, 1);
        fill_file(dentry, &syscall->rmdir.file);

        // the mount id of path_key is resolved by kprobe/mnt_want_write. It is already set by the time we reach this probe.
        key = syscall->rmdir.file.path_key;

        syscall->rmdir.dentry = dentry;
        syscall->policy = fetch_policy(EVENT_RMDIR);

        if (approve_syscall(syscall, rmdir_approvers) == DISCARDED) {
            // do not pop, we want to invalidate the inode even if the syscall is discarded
            return 0;
        }

        break;
    case EVENT_UNLINK:
        if (syscall->unlink.file.path_key.ino) {
            return 0;
        }

        // we resolve all the information before the file is actually removed
        dentry = (struct dentry *)CTX_PARM2(ctx);
        set_file_inode(dentry, &syscall->unlink.file, 1);
        fill_file(dentry, &syscall->unlink.file);

        // the mount id of path_key is resolved by kprobe/mnt_want_write. It is already set by the time we reach this probe.
        key = syscall->unlink.file.path_key;

        syscall->unlink.dentry = dentry;

        // fake rmdir event as we will generate and rmdir event at the end
        syscall->policy = fetch_policy(EVENT_RMDIR);

        if (approve_syscall(syscall, rmdir_approvers) == DISCARDED) {
            // do not pop, we want to invalidate the inode even if the syscall is discarded
            return 0;
        }

        break;
    default:
        return 0;
    }

    if (dentry != NULL) {
        syscall->resolver.key = key;
        syscall->resolver.dentry = dentry;
        syscall->resolver.discarder_event_type = dentry_resolver_discarder_event_type(syscall);
        syscall->resolver.callback = DR_SECURITY_INODE_RMDIR_CALLBACK_KPROBE_KEY;
        syscall->resolver.iteration = 0;
        syscall->resolver.ret = 0;

        resolve_dentry(ctx, KPROBE_OR_FENTRY_TYPE);

        // if the tail call fails, we need to pop the syscall cache entry
        pop_syscall_with(rmdir_predicate);
    }
    return 0;
}

TAIL_CALL_FNC(dr_security_inode_rmdir_callback, ctx_t *ctx) {
    struct syscall_cache_t *syscall = peek_syscall_with(rmdir_predicate);
    if (!syscall) {
        return 0;
    }

    if (syscall->resolver.ret == DENTRY_DISCARDED) {
        monitor_discarded(syscall->type);
        // do not pop, we want to invalidate the inode even if the syscall is discarded
        syscall->state = DISCARDED;
    }
    return 0;
}

int __attribute__((always_inline)) sys_rmdir_ret(void *ctx, int retval) {
    struct syscall_cache_t *syscall = pop_syscall_with(rmdir_predicate);
    if (!syscall) {
        return 0;
    }

    if (IS_UNHANDLED_ERROR(retval)) {
        return 0;
    }

    if (syscall->state != DISCARDED && is_event_enabled(EVENT_RMDIR)) {
        struct rmdir_event_t event = {
            .syscall.retval = retval,
            .syscall_ctx.id = syscall->ctx_id,
            .event.flags = syscall->async ? EVENT_FLAGS_ASYNC : 0,
            .file = syscall->rmdir.file,
        };

        struct proc_cache_t *entry = fill_process_context(&event.process);
        fill_container_context(entry, &event.container);
        fill_span_context(&event.span);

        send_event(ctx, EVENT_RMDIR, event);
    }

    if (retval >= 0) {
        expire_inode_discarders(syscall->rmdir.file.path_key.mount_id, syscall->rmdir.file.path_key.ino);
    }

    return 0;
}

HOOK_EXIT("do_rmdir")
int rethook_do_rmdir(ctx_t *ctx) {
    int retval = CTX_PARMRET(ctx);
    return sys_rmdir_ret(ctx, retval);
}

HOOK_SYSCALL_EXIT(rmdir) {
    int retval = SYSCALL_PARMRET(ctx);
    return sys_rmdir_ret(ctx, retval);
}

TAIL_CALL_TRACEPOINT_FNC(handle_sys_rmdir_exit, struct tracepoint_raw_syscalls_sys_exit_t *args) {
    return sys_rmdir_ret(args, args->ret);
}

#endif
