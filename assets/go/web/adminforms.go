package web

// Importing necessary packages
import (
	"fmt"
	"forum/assets/go/database"
	"net/http"
	"os"
)

func BanForm(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		// If there is an error, return an internal server error response
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	// Get the username, email, and password from the form data
	id := r.Form.Get("userId")
	username := r.Form.Get("username")
	banstatus := r.Form.Get("banstatus")

	if username != "QuentinDV" && username != "OwandjiD" && username != "Guest" {
		if banstatus == "true" {
			db, err := database.ConnectUserDB("db/database.db")
			if err != nil {
				return
			}
			defer db.Close()
			database.UnBanAccount(db, id)

		} else {
			db, err := database.ConnectUserDB("db/database.db")
			if err != nil {
				return
			}
			defer db.Close()
			database.BanAccount(db, id)
		}
	}

	// Redirect to the home page
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func ModeratorForm(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		// If there is an error, return an internal server error response
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	// Get the username, email, and password from the form data
	id := r.Form.Get("userId")
	username := r.Form.Get("username")
	moderator := r.Form.Get("moderator")

	if username != "QuentinDV" && username != "OwandjiD" && username != "Guest" {
		if moderator == "true" {
			db, err := database.ConnectUserDB("db/database.db")
			if err != nil {
				return
			}
			defer db.Close()
			database.DemoteFromModerator(db, id)

		} else {
			db, err := database.ConnectUserDB("db/database.db")
			if err != nil {
				return
			}
			defer db.Close()
			database.PromoteToModerator(db, id)
		}
	}

	// Redirect to the home page
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func AdminForm(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		// If there is an error, return an internal server error response
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	// Get the username, email, and password from the form data
	id := r.Form.Get("userId")
	admin := r.Form.Get("admin")
	username := r.Form.Get("username")

	if username != "QuentinDV" && username != "OwandjiD" && username != "Guest" {
		if admin == "true" {
			db, err := database.ConnectUserDB("db/database.db")
			if err != nil {
				return
			}
			defer db.Close()
			database.DemoteFromAdmin(db, id)
		} else {
			db, err := database.ConnectUserDB("db/database.db")
			if err != nil {
				return
			}
			defer db.Close()
			database.PromoteToAdmin(db, id)
		}
	}

	// Redirect to the home page
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func DeleteAccountForm(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Form data parsing error", http.StatusInternalServerError)
		return
	}

	// Get the username and user ID from the form data
	id := r.Form.Get("userId")
	username := r.Form.Get("username")

	// Only proceed if the user is not one of the protected usernames
	if username != "QuentinDV" && username != "OwandjiD" && username != "Guest" {
		db, err := database.ConnectUserDB("db/database.db")
		if err != nil {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// Get and delete the user's comments
		comments, err := database.GetCommentsByAccount(db, id)
		if err != nil {
			http.Error(w, "Error retrieving comments", http.StatusInternalServerError)
			return
		}

		// Get and delete the user's posts
		posts, err := database.GetPostsByCreator(db, id)
		if err != nil {
			http.Error(w, "Error retrieving posts", http.StatusInternalServerError)
			return
		}

		for _, post := range posts {
			err = database.DeletePost(db, post.PostID)
			if err != nil {
				http.Error(w, "Error deleting post", http.StatusInternalServerError)
				return
			}

			allAcc, err := database.GetAllAccounts(db)
			if err != nil {
				fmt.Println("Error getting all accounts:", err)
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}

			for _, acc := range allAcc {
				err = database.RemoveLikedPost(db, acc.Id, post.PostID)
				if err != nil {
					fmt.Println("Error removing liked post:", err)
					http.Error(w, "Database error", http.StatusInternalServerError)
					return
				}

				err = database.RemoveDisLikedPost(db, acc.Id, post.PostID)
				if err != nil {
					fmt.Println("Error removing disliked post:", err)
					http.Error(w, "Database error", http.StatusInternalServerError)
					return
				}

				err = database.RemoveSavedPost(db, acc.Id, post.PostID)
				if err != nil {
					fmt.Println("Error removing saved post:", err)
					http.Error(w, "Database error", http.StatusInternalServerError)
					return
				}

				for _, comment := range comments {
					err = database.RemoveLikedComment(db, acc.Id, comment.CommentID)
					if err != nil {
						fmt.Println("Error removing liked comment:", err)
						http.Error(w, "Database error", http.StatusInternalServerError)
						return
					}

					err = database.RemoveDislikedComment(db, acc.Id, comment.CommentID)
					if err != nil {
						fmt.Println("Error removing disliked comment:", err)
						http.Error(w, "Database error", http.StatusInternalServerError)
						return
					}
				}
			}

			// Remove the post image if it exists
			if post.ImageUrl != "" {
				fmt.Println("Deleting post image: " + post.ImageUrl)
				err = os.Remove(post.ImageUrl)
				if err != nil {
					http.Error(w, "Error deleting post image", http.StatusInternalServerError)
					return
				}
			}
		}

		// Delete the account
		err = database.DeleteAccount(db, id)
		if err != nil {
			http.Error(w, "Error deleting account", http.StatusInternalServerError)
			return
		}

		// Finally, delete the user's data
		err = database.DeleteUserData(db, id)
		if err != nil {
			http.Error(w, "Error deleting user data", http.StatusInternalServerError)
			return
		}

		// Redirect to the admin page
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	} else {
		http.Error(w, "Cannot delete protected user", http.StatusForbidden)
	}
}
