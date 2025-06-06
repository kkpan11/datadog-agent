// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build linux_bpf

package sm

import (
	"encoding/binary"
	"fmt"
)

func makeInstruction(op Op) codeFragment {
	switch op := op.(type) {
	case CallOp:
		return callInstruction{target: op.FunctionID}

	case ReturnOp:
		return staticInstruction{
			name:  "SM_OP_RETURN",
			bytes: []byte{},
		}

	case IllegalOp:
		return staticInstruction{
			name:  "SM_OP_ILLEGAL",
			bytes: []byte{},
		}

	case IncrementOutputOffsetOp:
		return staticInstruction{
			name:  "SM_OP_INCREMENT_OUTPUT_OFFSET",
			bytes: binary.LittleEndian.AppendUint32(nil, op.Value),
		}

	case ExprPrepareOp:
		return staticInstruction{
			name:  "SM_OP_EXPR_PREPARE",
			bytes: []byte{},
		}

	case ExprSaveOp:
		bytes := make([]byte, 0, 12)
		// Result offset and length.
		e := op.EventRootType.Expressions[op.ExprIdx]
		bytes = binary.LittleEndian.AppendUint32(bytes, e.Offset)
		bytes = binary.LittleEndian.AppendUint32(bytes, e.Expression.Type.GetByteSize())
		// Presence bit index.
		bytes = binary.LittleEndian.AppendUint32(bytes, op.ExprIdx)
		return staticInstruction{
			name:  "SM_OP_EXPR_SAVE",
			bytes: bytes,
		}

	case ExprDereferenceCfaOp:
		bytes := make([]byte, 0, 8)
		bytes = binary.LittleEndian.AppendUint32(bytes, op.Offset)
		bytes = binary.LittleEndian.AppendUint32(bytes, op.Len)
		return staticInstruction{
			name:  "SM_OP_EXPR_DEREFERENCE_CFA",
			bytes: bytes,
		}

	case ExprReadRegisterOp:
		return staticInstruction{
			name:  "SM_OP_EXPR_READ_REGISTER",
			bytes: []byte{op.Register, op.Size},
		}

	case ExprDereferencePtrOp:
		bytes := make([]byte, 0, 8)
		bytes = binary.LittleEndian.AppendUint32(bytes, op.Bias)
		bytes = binary.LittleEndian.AppendUint32(bytes, op.Len)
		return staticInstruction{
			name:  "SM_OP_EXPR_DEREFERENCE_PTR",
			bytes: bytes,
		}

	case ProcessPointerOp:
		return staticInstruction{
			name:  "SM_OP_PROCESS_POINTER",
			bytes: binary.LittleEndian.AppendUint32(nil, uint32(op.Pointee.GetID())),
		}

	case ProcessArrayPrepOp:
		return staticInstruction{
			name:  "SM_OP_PROCESS_ARRAY_PREP",
			bytes: binary.LittleEndian.AppendUint32(nil, op.Array.Count),
		}

	case ProcessArrayRepeatOp:
		return staticInstruction{
			name:  "SM_OP_PROCESS_ARRAY_PREP",
			bytes: binary.LittleEndian.AppendUint32(nil, op.OffsetShift),
		}

	case ProcessSliceOp:
		return staticInstruction{
			name:  "SM_OP_PROCESS_SLICE",
			bytes: binary.LittleEndian.AppendUint32(nil, uint32(op.SliceData.GetID())),
		}

	case ProcessSliceDataPrepOp:
		return staticInstruction{
			name:  "SM_OP_PROCESS_SLICE_DATA_PREP",
			bytes: []byte{},
		}

	case ProcessSliceDataRepeatOp:
		return staticInstruction{
			name:  "SM_OP_PROCESS_SLICE_DATA_REPEAT",
			bytes: binary.LittleEndian.AppendUint32(nil, op.OffsetShift),
		}

	case ProcessStringOp:
		return staticInstruction{
			name:  "SM_OP_PROCESS_STRING",
			bytes: []byte{},
		}

	case ProcessGoEmptyInterfaceOp:
		return staticInstruction{
			name:  "SM_OP_PROCESS_GO_EMPTY_INTERFACE",
			bytes: []byte{},
		}

	case ProcessGoInterfaceOp:
		return staticInstruction{
			name:  "SM_OP_PROCESS_GO_INTERFACE",
			bytes: []byte{},
		}

	case ProcessGoHmapOp:
		return staticInstruction{
			name:  "SM_OP_PROCESS_GO_HMAP",
			bytes: binary.LittleEndian.AppendUint32(nil, uint32(op.BucketsArray.GetID())),
		}

	case ProcessGoSwissMapOp:
		bytes := make([]byte, 0, 8)
		bytes = binary.LittleEndian.AppendUint32(bytes, uint32(op.TablePtrSlice.GetID()))
		bytes = binary.LittleEndian.AppendUint32(bytes, uint32(op.Group.GetID()))
		return staticInstruction{
			name:  "SM_OP_PROCESS_GO_SWISS_MAP",
			bytes: bytes,
		}

	case ProcessGoSwissMapGroupsOp:
		return staticInstruction{
			name:  "SM_OP_PROCESS_GO_SWISS_MAP_GROUPS",
			bytes: binary.LittleEndian.AppendUint32(nil, uint32(op.Group.GetID())),
		}

	case ChasePointersOp:
		return staticInstruction{
			name:  "SM_OP_CHASE_POINTERS",
			bytes: []byte{},
		}

	case PrepareEventRootOp:
		bytes := make([]byte, 0, 8)
		bytes = binary.LittleEndian.AppendUint32(bytes, uint32(op.EventRootType.GetID()))
		bytes = binary.LittleEndian.AppendUint32(bytes, op.EventRootType.GetByteSize())
		return staticInstruction{
			name:  "SM_OP_PREPARE_EVENT_ROOT",
			bytes: bytes,
		}

	default:
		panic(fmt.Sprintf("unsupported op: %T", op))
	}
}
