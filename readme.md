    DisplayName string  `json:"display_name,omitempty" gorm:"column:display_name;<-:false"`


	Phone        *string `json:"phone"` // 因为零值为 “”， 会触发DB的phone unique constraint Error， 所以这里使用指针实现可有可无


    ERROR: cached plan must not change result type (SQLSTATE 0A000)
    