package helper

import (
	"fmt"
	"github.com/devtron-labs/devtron/api/bean"
	bean2 "github.com/devtron-labs/devtron/pkg/auth/user/bean"
	bean3 "github.com/devtron-labs/devtron/pkg/timeoutWindow/repository/bean"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type UserRepositoryQueryBuilder struct {
	logger *zap.SugaredLogger
}

func NewUserRepositoryQueryBuilder(logger *zap.SugaredLogger) UserRepositoryQueryBuilder {
	userListingRepositoryQueryBuilder := UserRepositoryQueryBuilder{
		logger: logger,
	}
	return userListingRepositoryQueryBuilder
}

const (
	QueryTimeFormat      string = "2006-01-02 15:04:05-07:00"
	TimeStampFormat      string = "YYYY-MM-DD HH24:MI:SS"
	TimeFormatForParsing string = "2006-01-02 15:04:05 -0700 MST"
)

func (impl UserRepositoryQueryBuilder) GetQueryForUserListingWithFilters(req *bean.FetchListingRequest) string {
	whereCondition := fmt.Sprintf("where active = %t AND (user_type is NULL or user_type != '%s') ", true, bean.USER_TYPE_API_TOKEN)
	orderCondition := ""
	whereCondition += impl.buildQueryForStatusFilter(req.StatusType, req.CurrentTime)
	if len(req.SearchKey) > 0 {
		emailIdLike := "%" + req.SearchKey + "%"
		whereCondition += fmt.Sprintf("AND email_id ilike '%s' ", emailIdLike)
	}

	if len(req.SortBy) > 0 && !req.CountCheck {
		orderCondition += fmt.Sprintf("order by %s ", req.SortBy)
		if req.SortOrder == bean2.Desc {
			orderCondition += string(req.SortOrder)
		}
	}

	if req.Size > 0 && !req.CountCheck {
		orderCondition += " limit " + strconv.Itoa(req.Size) + " offset " + strconv.Itoa(req.Offset) + ""
	}
	var query string
	if req.CountCheck {
		query = fmt.Sprintf("select count(*) from users AS user_model left join user_audit AS au on au.user_id=user_model.id left join timeout_window_configuration AS timeout_window_configuration on timeout_window_configuration.id=user_model.timeout_window_configuration_id %s %s;", whereCondition, orderCondition)
	} else {
		// have not collected client ip here. always will be empty
		query = fmt.Sprintf(`SELECT "user_model".*, "timeout_window_configuration"."id" AS "timeout_window_configuration__id", "timeout_window_configuration"."timeout_window_expression" AS "timeout_window_configuration__timeout_window_expression", "timeout_window_configuration"."timeout_window_expression_format" AS "timeout_window_configuration__timeout_window_expression_format", "user_audit"."id" AS "user_audit__id", "user_audit"."updated_on" AS "user_audit__updated_on","user_audit"."user_id" AS "user_audit__user_id" ,"user_audit"."created_on" AS "user_audit__created_on" ,"user_audit"."updated_on" AS "last_login" from users As "user_model" LEFT JOIN user_audit As "user_audit" on "user_audit"."user_id" = "user_model"."id" LEFT JOIN timeout_window_configuration AS "timeout_window_configuration" ON "timeout_window_configuration"."id" = "user_model"."timeout_window_configuration_id" %s %s;`, whereCondition, orderCondition)
	}

	return query
}

func (impl UserRepositoryQueryBuilder) buildQueryForStatusFilter(statusType bean2.StatusType, currentTime time.Time) string {
	condition := ""
	//formatted for query comparison
	formattedTimeForQuery := currentTime.Format(QueryTimeFormat)

	if statusType == bean2.Active {
		condition += "AND (user_model.timeout_window_configuration_id is null ) "
	} else if statusType == bean2.Inactive {
		condition += fmt.Sprintf("AND (timeout_window_configuration.timeout_window_expression_format = %v AND timeout_window_configuration.timeout_window_expression < '%s' ) ", bean3.TimeZeroFormat, formattedTimeForQuery)
	} else if statusType == bean2.TemporaryAccess {
		condition += fmt.Sprintf(" AND (timeout_window_configuration.timeout_window_expression_format = %v AND timeout_window_configuration.timeout_window_expression > '%s' ) ", bean3.TimeStamp, formattedTimeForQuery)
	} else if statusType == bean2.Active_TemporaryAccess {
		condition += fmt.Sprintf("AND (user_model.timeout_window_configuration_id is null OR ( timeout_window_configuration.timeout_window_expression_format = %v AND timeout_window_configuration.timeout_window_expression > '%s' ) ) ", bean3.TimeStamp, formattedTimeForQuery)
	} else if statusType == bean2.Active_InActive {
		condition += fmt.Sprintf("AND (user_model.timeout_window_configuration_id is null OR (timeout_window_configuration.timeout_window_expression_format = %v AND timeout_window_configuration.timeout_window_expression < '%s') ) ", bean3.TimeZeroFormat, formattedTimeForQuery)
	} else if statusType == bean2.Inactive_TemporaryAccess {
		condition += fmt.Sprintf("AND ((timeout_window_configuration.timeout_window_expression_format = %v AND timeout_window_configuration.timeout_window_expression < '%s') OR ( timeout_window_configuration.timeout_window_expression_format = %v AND timeout_window_configuration.timeout_window_expression > '%s' ) ) ", bean3.TimeZeroFormat, formattedTimeForQuery, bean3.TimeStamp, formattedTimeForQuery)
	}
	return condition
}

func (impl UserRepositoryQueryBuilder) GetQueryForAllUserWithAudit() string {
	whereCondition := fmt.Sprintf("where active = %t AND (user_type is NULL or user_type != '%s') ", true, bean.USER_TYPE_API_TOKEN)
	orderCondition := fmt.Sprintf("order by user_model.updated_on %s", bean2.Desc)
	query := fmt.Sprintf(`SELECT "user_model".*, "timeout_window_configuration"."id" AS "timeout_window_configuration__id", "timeout_window_configuration"."timeout_window_expression" AS "timeout_window_configuration__timeout_window_expression", "timeout_window_configuration"."timeout_window_expression_format" AS "timeout_window_configuration__timeout_window_expression_format", "user_audit"."id" AS "user_audit__id", "user_audit"."updated_on" AS "user_audit__updated_on","user_audit"."user_id" AS "user_audit__user_id" ,"user_audit"."created_on" AS "user_audit__created_on" from users As "user_model" LEFT JOIN user_audit As "user_audit" on "user_audit"."user_id" = "user_model"."id" LEFT JOIN timeout_window_configuration AS "timeout_window_configuration" ON "timeout_window_configuration"."id" = "user_model"."timeout_window_configuration_id" %s %s;`, whereCondition, orderCondition)
	return query
}

func (impl UserRepositoryQueryBuilder) GetQueryForGroupListingWithFilters(req *bean.FetchListingRequest) string {
	whereCondition := fmt.Sprintf("where active = %t ", true)
	orderCondition := ""
	if len(req.SearchKey) > 0 {
		nameIdLike := "%" + req.SearchKey + "%"
		whereCondition += fmt.Sprintf("AND name ilike '%s' ", nameIdLike)
	}

	if len(req.SortBy) > 0 && !req.CountCheck {
		orderCondition += fmt.Sprintf("order by %s ", req.SortBy)
		if req.SortOrder == bean2.Desc {
			orderCondition += string(req.SortOrder)
		}
	}

	if req.Size > 0 && !req.CountCheck {
		orderCondition += " limit " + strconv.Itoa(req.Size) + " offset " + strconv.Itoa(req.Offset) + ""
	}
	var query string
	if req.CountCheck {
		query = fmt.Sprintf("SELECT count(*) from role_group %s %s;", whereCondition, orderCondition)
	} else {
		query = fmt.Sprintf("SELECT * from role_group %s %s;", whereCondition, orderCondition)
	}
	return query

}
