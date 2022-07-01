package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/shijuvar/gokit/src/testing/httpbdd/controllers"
	"github.com/shijuvar/gokit/src/testing/httpbdd/model"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserController", func() {
	var (
		r       *mux.Router
		w       *httptest.ResponseRecorder
		store   *FakeUserStore
		handler controllers.Handler
	)
	BeforeEach(func() {
		r = mux.NewRouter()
		store = newFakeUserStore()
		handler = controllers.Handler{
			// Injecting a test stub
			// In production code, this would be a persistent store
			Store: store,
		}
	})

	// Specs for HTTP Get to "/users"
	Describe("Get list of Users", func() {
		Context("Get all Users from data store", func() {
			It("Should get list of Users", func() {
				r.HandleFunc("/users", handler.GetUsers).Methods("GET")
				req, err := http.NewRequest("GET", "/users", nil)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				r.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))
				var users []model.User
				json.Unmarshal(w.Body.Bytes(), &users)
				// Verifying mocked data of 2 users
				Expect(len(users)).To(Equal(2))
			})
		})
	})

	// Specs for HTTP Post to "/users"
	Describe("Post a new User", func() {
		Context("Provide a valid User data", func() {
			It("Should create a new User and get HTTP Status: 201", func() {
				r.HandleFunc("/users", handler.CreateUser).Methods("POST")
				userJson := `{"firstname": "Alex", "lastname": "John", "email": "alex@xyz.com"}`

				req, err := http.NewRequest(
					"POST",
					"/users",
					strings.NewReader(userJson),
				)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				r.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusCreated))
			})
		})
		Context("Provide a User data that contains duplicate email id", func() {
			It("Should get HTTP Status: 400", func() {
				r.HandleFunc("/users", handler.CreateUser).Methods("POST")
				userJson := `{"firstname": "Shiju", "lastname": "Varghese", "email": "shiju@xyz.com"}`

				req, err := http.NewRequest(
					"POST",
					"/users",
					strings.NewReader(userJson),
				)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				r.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})
})

// FakeUserStore provides a mocked implementation of interface model.UserStore
type FakeUserStore struct {
	userStore []model.User
}

// GetUsers returns all users
func (store *FakeUserStore) GetUsers() []model.User {
	return store.userStore
}

// AddUser inserts a User
func (store *FakeUserStore) AddUser(user model.User) error {
	// Check whether email is exists
	for _, u := range store.userStore {
		if u.Email == user.Email {
			return model.ErrorEmailExists
		}
	}
	store.userStore = append(store.userStore, user)
	return nil
}

// newFakeUserStore provides two dummy data for Users
func newFakeUserStore() *FakeUserStore {
	store := &FakeUserStore{}
	store.AddUser(model.User{
		FirstName: "Shiju",
		LastName:  "Varghese",
		Email:     "shiju@xyz.com",
	})

	store.AddUser(model.User{
		FirstName: "Irene",
		LastName:  "Rose",
		Email:     "irene@xyz.com",
	})
	return store
}
