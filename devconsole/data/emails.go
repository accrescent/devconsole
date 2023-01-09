package data

import (
	"context"
	"regexp"

	"github.com/google/go-github/v49/github"
)

// GitHub's noreply email format is documented at
// https://docs.github.com/en/account-and-profile/setting-up-and-managing-your-personal-account-on-github/managing-email-preferences/setting-your-commit-email-address
var noReplyEmail = regexp.MustCompile(`^([0-9]{7}\+)?.*@users\.noreply\.github\.com$`)

func GetUsableEmails(ctx context.Context, client *github.Client) ([]string, error) {
	var usableEmails []string

	emails, _, err := client.Users.ListEmails(ctx, nil)
	if err != nil {
		return []string{}, err
	}
	for _, email := range emails {
		address := email.GetEmail()
		if email.GetVerified() && !noReplyEmail.MatchString(address) {
			usableEmails = append(usableEmails, address)
		}
	}

	return usableEmails, nil
}
