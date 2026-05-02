package shared

import (
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

const RootNestedName = "root"

func ValidateMask(mask *fieldmaskpb.FieldMask, msg proto.Message, immutableFields map[string]bool) error {
	if mask == nil || len(mask.GetPaths()) == 0 {
		return status.Error(codes.InvalidArgument, "update_mask is required")
	}
	if msg == nil {
		return status.Error(codes.InvalidArgument, "message is required")
	}

	for _, path := range mask.GetPaths() {
		fd, ok := descriptorFieldByPath(msg.ProtoReflect().Descriptor(), path)
		if !ok {
			return status.Errorf(codes.InvalidArgument, "unknown field: %s", path)
		}
		if isMessage(fd) {
			return status.Errorf(codes.InvalidArgument, "field %s is not updatable", path)
		}
		if immutableFields[path] {
			return status.Errorf(codes.InvalidArgument, "field %s is not updatable", path)
		}
	}

	return nil
}

func ExtractChanges(mask *fieldmaskpb.FieldMask, msg proto.Message) map[string]any {
	changes := make(map[string]any, len(mask.GetPaths()))
	ref := msg.ProtoReflect()

	for _, path := range mask.GetPaths() {
		changes[path] = valueByPath(ref, path)
	}

	return changes
}

func ExtractNestedChanges(changes map[string]any, fields map[string]string, nestedName string) map[string]any {
	out := make(map[string]any)

	for path, value := range changes {
		field, ok := fieldNameForNestedPath(path, nestedName)
		if !ok {
			continue
		}

		target, ok := fields[field]
		if !ok {
			continue
		}

		out[target] = value
	}

	return out
}

func fieldNameForNestedPath(path string, nestedName string) (string, bool) {
	if path == "" {
		return "", false
	}

	if nestedName == "" || nestedName == RootNestedName {
		if strings.Contains(path, ".") {
			return "", false
		}
		return path, true
	}

	prefix := nestedName + "."
	if !strings.HasPrefix(path, prefix) {
		return "", false
	}

	field := strings.TrimPrefix(path, prefix)
	if field == "" || strings.Contains(field, ".") {
		return "", false
	}

	return field, true
}

func descriptorFieldByPath(root protoreflect.MessageDescriptor, path string) (protoreflect.FieldDescriptor, bool) {
	if path == "" {
		return nil, false
	}

	descriptor := root
	parts := strings.Split(path, ".")
	for i, part := range parts {
		fd := descriptor.Fields().ByName(protoreflect.Name(part))
		if fd == nil {
			return nil, false
		}

		if i == len(parts)-1 {
			return fd, true
		}

		if !isMessage(fd) {
			return nil, false
		}
		descriptor = fd.Message()
	}

	return nil, false
}

func valueByPath(root protoreflect.Message, path string) any {
	if path == "" {
		return nil
	}

	ref := root
	parts := strings.Split(path, ".")
	for i, part := range parts {
		fd := ref.Descriptor().Fields().ByName(protoreflect.Name(part))
		if fd == nil {
			return nil
		}

		if i == len(parts)-1 {
			if fd.HasPresence() && !ref.Has(fd) {
				return nil
			}
			return ref.Get(fd).Interface()
		}

		if !isMessage(fd) || !ref.Has(fd) {
			return nil
		}
		ref = ref.Get(fd).Message()
	}

	return nil
}

func isMessage(fd protoreflect.FieldDescriptor) bool {
	return fd.Kind() == protoreflect.MessageKind || fd.Kind() == protoreflect.GroupKind
}
