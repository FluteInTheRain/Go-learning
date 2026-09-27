package main

import (
	"time"
)

func (s *store) Create(p Post) Post {
	s.mu.Lock()
	defer s.mu.Unlock()

	p.ID = s.nextPostID
	p.CreatedAt = time.Now()
	s.posts[p.ID] = p
	s.nextPostID++
	return p
}

func (s *store) List() []Post {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Post, 0, len(s.posts))
	for _, post := range s.posts {
		result = append(result, post)
	}
	return result
}

func (s *store) Get(id int) (Post, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.posts[id]
	return p, ok
}

func (s *store) Update(id int, p Post) (Post, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.posts[id]
	if !ok {
		return Post{}, false
	}
	existing.Title = p.Title
	existing.Body = p.Body
	s.posts[id] = existing
	return existing, true
}

// Delete removes a post; returns false when it does not exist.
func (s *store) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.posts[id]; !ok {
		return false
	}
	delete(s.posts, id)
	return true
}

// CreateUser adds a new user and returns it with ID assigned.
func (s *store) CreateUser(email, passwordHash, role string) User {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := User{
		ID:           s.nextUserID,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    time.Now(),
	}
	s.users[user.ID] = user
	s.nextUserID++
	return user
}

// FindUserByEmail looks up a user by email; returns nil when not found.
func (s *store) FindUserByEmail(email string) *User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.Email == email {
			return &user
		}
	}
	return nil
}
