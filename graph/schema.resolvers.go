package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.55

import (
	"context"
	"fmt"
	"gogqlauth/graph/database"
	"gogqlauth/graph/middlewares"
	"gogqlauth/graph/model"
	"gogqlauth/graph/utils"
	"strings"

	surrealdb "github.com/surrealdb/surrealdb.go"
)

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, input model.NewUser) (*model.User, error) {
	result, err := database.DB.Query(`
    SELECT * FROM user WHERE username=$username OR email=$email;`, map[string]interface{}{
		"username": input.Username,
		"email":    strings.ToLower(input.Email),
	})
	if err != nil {
		return nil, err
	}
	users, err := surrealdb.SmartUnmarshal[[]model.User](result, nil)
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return nil, fmt.Errorf("email and username should be unique")
	}
	result, err = database.DB.Query(
		`CREATE ONLY user:ulid() 
		SET username=$username,
		email=$email,
		password=crypto::argon2::generate($password),
		createdAt=time::now(),
		updatedAt=time::now();`, map[string]interface{}{
			"username": input.Username,
			"email":    strings.ToLower(input.Email),
			"password": input.Password,
		})
	if err != nil {
		return nil, err
	}

	newUser, err := surrealdb.SmartUnmarshal[model.User](result, nil)
	if err != nil {
		return nil, err
	}
	return &newUser, nil
}

// LoginWithEmailPassword is the resolver for the loginWithEmailPassword field.
func (r *mutationResolver) LoginWithEmailPassword(ctx context.Context, email string, password string, token string) (*model.Token, error) {
	// Query the database for a user with matching email and password
	result, err := database.DB.Query(`
    SELECT * FROM user WHERE email=$email AND crypto::argon2::compare(password, $pass);`, map[string]interface{}{
		"email": strings.ToLower(email),
		"pass":  password,
	})
	if err != nil {
		return nil, err
	}
	users, err := surrealdb.SmartUnmarshal[[]model.User](result, nil)
	if err != nil {
		return nil, err
	}
	if len(users) <= 0 {
		return nil, fmt.Errorf("user not found")
	}
	user := users[0]
	if user.Password == "" {
		return nil, fmt.Errorf("user not available for password sign in, use other sign in methods instead")
	}
	// Generate and return a new token
	mytoken, err := utils.HandleLogin(&user, ctx)
	if err != nil {
		return nil, err
	}
	return mytoken, nil
}

// ChangePassword is the resolver for the changePassword field.
func (r *mutationResolver) ChangePassword(ctx context.Context, token string, newPassword string) (bool, error) {
	// Retrieve password change request using the token
	data, err := database.DB.Select("passwordChange:" + token)
	if err != nil {
		return false, nil
	}
	passwordChange, err := surrealdb.SmartUnmarshal[model.PasswordChange](data, nil)
	if err != nil {
		return false, nil
	}
	// Update the user's password in the database
	_, err = database.DB.Query("UPDATE $id SET password=crypto::argon2::generate($password);", map[string]interface{}{
		"id":       passwordChange.User,
		"password": newPassword,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

// RequestChangePassword is the resolver for the requestChangePassword field.
func (r *mutationResolver) RequestChangePassword(ctx context.Context, token string, email string) (bool, error) {
	// Find the user by email
	result, err := database.DB.Query(`
    SELECT * FROM user WHERE email=$email;`, map[string]interface{}{
		"email": email,
	})
	if err != nil {
		return false, err
	}
	users, err := surrealdb.SmartUnmarshal[[]model.User](result, nil)
	if err != nil {
		return false, err
	}
	if len(users) <= 0 {
		return false, fmt.Errorf("user not found")
	}
	// Create a new password change request
	result2, err := database.DB.Query(
		`CREATE ONLY passwordChange:ulid() SET userId=$userId;`, map[string]interface{}{
			"userId": users[0].ID,
		})
	if err != nil {
		return false, err
	}
	results2, err := surrealdb.SmartUnmarshal[model.PasswordChange](result2, nil)
	if err != nil {
		return false, err
	}
	// Send password reset email
	utils.SendPasswordResetEmail(users[0].Email, users[0].Username, strings.Split(results2.ID, ":")[1])
	return true, nil
}

// RefreshToken is the resolver for the refreshToken field.
func (r *mutationResolver) RefreshToken(ctx context.Context, accessToken string, refreshToken string, device string) (*model.Token, error) {
	// Find the existing token in the database
	result, err := database.DB.Query("SELECT * FROM token WHERE accessToken=$accessToken AND refreshToken=$refreshToken;", map[string]string{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
	if err != nil {
		return nil, err
	}

	results, err := surrealdb.SmartUnmarshal[[]model.Token](result, nil)
	if err != nil {
		return nil, err
	}
	if len(results) <= 0 {
		return nil, fmt.Errorf("token not found")
	}
	token := results[0]
	// Retrieve the user associated with the token
	var user model.User
	data, err := database.DB.Select(token.User)
	if err != nil {
		return nil, err
	}
	user, err = surrealdb.SmartUnmarshal[model.User](data, nil)
	if err != nil {
		return nil, err
	}
	// Generate a new token
	newToken, err := utils.HandleLogin(&user, ctx)
	if err != nil {
		return nil, err
	}
	// Delete the old token
	database.DB.Delete(token.ID)
	return newToken, nil
}

// SomeMethod is the resolver for the someMethod field.
func (r *mutationResolver) SomeMethod(ctx context.Context, input string) (string, error) {
	// Check if the user is authenticated
	user := middlewares.ForContext(ctx)
	if user == nil {
		return "", fmt.Errorf("access denied")
	}
	// Do something with the authenticated user
	return "Hello, " + user.Username, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context) (*model.User, error) {
	// Check if the user is authenticated
	user := middlewares.ForContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("access denied")
	}
	return user, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
