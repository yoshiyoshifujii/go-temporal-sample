package saga

import (
	"context"
	"fmt"
)

type (
	ActivityInput struct {
		IsError             bool
		IsCompensationError bool
	}
)

func Activity01(ctx context.Context, input ActivityInput) error {
	if input.IsError {
		return fmt.Errorf("activity01 error")
	} else {
		fmt.Println("activity01 success")
		return nil
	}
}

func Activity01Compensation(ctx context.Context, input ActivityInput) error {
	if input.IsCompensationError {
		return fmt.Errorf("activity01 compensation error")
	} else {
		fmt.Println("activity01 compensation success")
		return nil
	}
}

func Activity02(ctx context.Context, input ActivityInput) error {
	if input.IsError {
		return fmt.Errorf("activity02 error")
	} else {
		fmt.Println("activity02 success")
		return nil
	}
}

func Activity02Compensation(ctx context.Context, input ActivityInput) error {
	if input.IsCompensationError {
		return fmt.Errorf("activity02 compensation error")
	} else {
		fmt.Println("activity02 compensation success")
		return nil
	}
}

func Activity03(ctx context.Context, input ActivityInput) error {
	if input.IsError {
		return fmt.Errorf("activity03 error")
	} else {
		fmt.Println("activity03 success")
		return nil
	}
}
