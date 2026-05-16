package user

import (
	"context"

	commonv1 "grpc-sandbox/gen/grpc/common/v1"
	userv1 "grpc-sandbox/gen/grpc/user/v1"
	"grpc-sandbox/internal/database"

	orm "github.com/kitti12911/lib-orm/v3"
	"github.com/kitti12911/lib-util/v3/fieldmask"
	"github.com/kitti12911/lib-util/v3/pagination"
)

var immutableUserFields = map[string]bool{
	"id":         true,
	"created_at": true,
	"updated_at": true,
}

type userService interface {
	GetByID(ctx context.Context, params GetByIDParams) (*database.User, error)
	Create(ctx context.Context, params CreateParams) (string, error)
	Update(ctx context.Context, params UpdateParams) (int64, error)
	Patch(ctx context.Context, params PatchParams) (int64, error)
	List(ctx context.Context, params ListParams) (*ListResult, error)
	Delete(ctx context.Context, params DeleteParams) (int64, error)
}

type Handler struct {
	userv1.UnimplementedUserServiceServer
	userService userService
}

func NewHandler(userService userService) *Handler {
	return &Handler{userService: userService}
}

func (h *Handler) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	result, err := h.userService.GetByID(ctx, GetByIDParams{
		ID: req.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &userv1.GetUserResponse{
		User: toProtoUser(result),
	}, nil
}

func (h *Handler) ListUsers(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
	pag := req.GetPagination()
	page := pagination.ParseInput(pag.GetPage(), pag.GetPageSize())

	result, err := h.userService.List(ctx, ListParams{
		Limit:   page.Limit,
		Offset:  page.Offset,
		Filters: orm.FiltersFromProto(req.GetFilters()),
		OrderBy: orm.OrderByFromProto(req.GetOrderBy()),
	})
	if err != nil {
		return nil, err
	}

	users := make([]*userv1.User, len(result.Users))
	for i := range result.Users {
		users[i] = toProtoUser(&result.Users[i])
	}

	pageOut := pagination.CalcOutput(pag.GetPage(), pag.GetPageSize(), result.Total)

	return &userv1.ListUsersResponse{
		Users: users,
		Pagination: &commonv1.PaginationResponse{
			Page:       pageOut.Page,
			PageSize:   pageOut.PageSize,
			TotalPages: pageOut.TotalPages,
			TotalSize:  pageOut.TotalSize,
		},
	}, nil
}

func (h *Handler) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	id, err := h.userService.Create(ctx, createParamsFromProto(req.GetUser()))
	if err != nil {
		return nil, err
	}

	return &userv1.CreateUserResponse{Id: id}, nil
}

func (h *Handler) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	affectedRows, err := h.userService.Update(ctx, updateParamsFromProto(req.GetId(), req.GetUser()))
	if err != nil {
		return nil, err
	}

	return &userv1.UpdateUserResponse{AffectedRows: affectedRows}, nil
}

func (h *Handler) PatchUser(ctx context.Context, req *userv1.PatchUserRequest) (*userv1.PatchUserResponse, error) {
	if err := fieldmask.ValidateMask(req.GetUpdateMask(), req.GetUser(), immutableUserFields); err != nil {
		return nil, err
	}

	affectedRows, err := h.userService.Patch(ctx, patchParamsFromProto(req.GetId(), req.GetUser(), req.GetUpdateMask()))
	if err != nil {
		return nil, err
	}

	return &userv1.PatchUserResponse{AffectedRows: affectedRows}, nil
}

func (h *Handler) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	affectedRows, err := h.userService.Delete(ctx, DeleteParams{
		ID: req.GetId(),
	})
	if err != nil {
		return nil, err
	}

	return &userv1.DeleteUserResponse{AffectedRows: affectedRows}, nil
}
