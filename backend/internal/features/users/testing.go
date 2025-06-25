package users

func GetTestUser() *SignInResponse {
	isAnyUserExists, err := userService.IsAnyUserExist()
	if err != nil {
		panic(err)
	}

	if !isAnyUserExists {
		err = userService.SignUp(&SignUpRequest{
			Email:    "test@test.com",
			Password: "test",
		})

		if err != nil {
			panic(err)
		}
	}

	user, err := userService.GetFirstUser()
	if err != nil {
		panic(err)
	}

	signInResponse, err := userService.GenerateAccessToken(user)
	if err != nil {
		panic(err)
	}

	return signInResponse
}
