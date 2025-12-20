INSERT INTO users (id, email, username, password) VALUES 
('user_01', 'reader@test.com', 'BookWorm99', 'hashedpass'),
('user_02', 'critic@test.com', 'SeriousCritic', 'hashedpass');

INSERT INTO posts (user_id, title, content, post_type, book_id, chapter, created_at) VALUES 
('user_01', 'Welcome to the club', 'Happy to be here discussing books!', 'discussion', NULL, NULL, NOW()),

('user_02', 'The Hobbit - A Classic', 'Ideally, this post reviews the book...', 'review', 101, NULL, NOW()),

('user_01', 'Deep dive into Chapter 3', 'The symbolism here was intense.', 'analysis', 101, '3', NOW());

INSERT INTO comments (post_id, user_id, content, created_at) VALUES 
(1, 'user_02', 'Glad to have you here!', NOW()),
(2, 'user_01', 'I totally agree with this review.', NOW());

INSERT INTO votes (user_id, target_type, target_id, value) VALUES 
('user_02', 'post', 1, 1),
('user_01', 'post', 2, 1),
('user_02', 'comment', 1, -1);