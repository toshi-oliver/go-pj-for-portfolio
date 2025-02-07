package validator

import (
	"go-pj-for-portfolio/model"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func TaskValidate(task model.Task) error {
	return validation.ValidateStruct(&task,
		validation.Field(
			&task.Title,
			validation.Required.Error("title is required"),
			validation.RuneLength(1, 10).Error("limited max 10 char"),
		),
	)
}
