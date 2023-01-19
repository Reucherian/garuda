package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"garuda.com/m/cmd/utils"
	"garuda.com/m/model"
	"google.golang.org/grpc"
)

var authClient model.AuthServiceClient
var activityClient model.ActivityServiceClient

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		_, err := authClient.Login(r.Context(), &model.UserRequest{Username: r.Form.Get("username"), Password: r.Form.Get("password")})
		if err != nil {
			fmt.Fprintf(w, "Login Failed : %s", err)
		} else {
			token, time, err := utils.GenerateJWT(r.Form.Get("username"))
			if err != nil {
				fmt.Fprintf(w, "Failed to generate token : %s", err)
			} else {
				http.SetCookie(w, &http.Cookie{Name: "token", Value: token, Expires: time})
				http.Redirect(w, r, "/home", http.StatusFound)
			}
		}
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("register.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		_, err := authClient.Register(r.Context(), &model.UserRequest{Username: r.Form.Get("username"), Password: r.Form.Get("password")})
		if err != nil {
			fmt.Fprintf(w, "Register Failed : %s", err)
			fmt.Println("error:", err)
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	}
}

type HomeContext struct {
	Username  string
	Posts     []map[string]string
	Following int
}

func home(w http.ResponseWriter, r *http.Request, claims *utils.Claims) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./home.gtpl")
		followers, err := activityClient.GetFollowings(r.Context(), &model.User{Username: claims.Username})
		if err != nil {
			fmt.Println("error:", err)
		}
		AllPosts := []map[string]string{}
		users := followers.GetUsers()
		for _, user := range users {
			posts, err := activityClient.GetPosts(r.Context(), &model.User{Username: user.Username})
			if err != nil {
				fmt.Println("error:", err)
			} else {
				for _, post := range posts.Posts {
					AllPosts = append(AllPosts, map[string]string{"author": user.Username, "title": post.Title, "content": post.Content})
				}
			}
		}
		context := HomeContext{Username: claims.Username, Posts: AllPosts, Following: len(users)}
		t.Execute(w, context)
	}
}

func createPost(w http.ResponseWriter, r *http.Request, claims *utils.Claims) {
	if r.Method == "POST" {
		r.ParseForm()
		post := model.Post{Title: r.Form.Get("title"), Content: r.Form.Get("content")}
		user := model.User{Username: claims.Username}
		_, err := activityClient.CreatePost(r.Context(), &model.PostRequest{User: &user, Post: &post})
		if err != nil {
			fmt.Println("error:", err)
		}
		http.Redirect(w, r, "/home", http.StatusFound)
	}
}

func followUser(w http.ResponseWriter, r *http.Request, claims *utils.Claims) {
	if r.Method == "POST" {
		r.ParseForm()
		user := model.User{Username: claims.Username}
		following := model.User{Username: r.Form.Get("username")}
		_, err := activityClient.AddFollowing(r.Context(), &model.FollowingRequest{User: &user, Following: &following})
		if err != nil {
			if err.Error() == "Cannot follow yourself" {
				fmt.Fprintf(w, "Cannot follow yourself")
				return
			}
			fmt.Println("error:", err)
		}
		http.Redirect(w, r, "/home", http.StatusFound)
	}
}

func deleteFollowing(w http.ResponseWriter, r *http.Request, claims *utils.Claims) {
	if r.Method == "POST" {
		r.ParseForm()
		user := model.User{Username: claims.Username}
		following := model.User{Username: r.Form.Get("username")}
		_, err := activityClient.DeleteFollowing(r.Context(), &model.FollowingRequest{User: &user, Following: &following})
		if err != nil {
			fmt.Println("error:", err)
		}
		http.Redirect(w, r, "/home", http.StatusFound)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Expires: time.Now(),
		})
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func main() {
	authConn, err := grpc.Dial("127.0.0.1:4040", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	activityConn, err := grpc.Dial("127.0.0.1:4043", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	authClient = model.NewAuthServiceClient(authConn)
	activityClient = model.NewActivityServiceClient(activityConn)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/home", utils.VerifyJWT(home))
	http.HandleFunc("/createPost", utils.VerifyJWT(createPost))
	http.HandleFunc("/followUser", utils.VerifyJWT(followUser))
	http.HandleFunc("/deleteFollowing", utils.VerifyJWT(deleteFollowing))
	http.HandleFunc("/logout", logout)
	err = http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
