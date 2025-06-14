USE `social-todo-db`;

-- Insert Accounts
INSERT INTO Account (Password, Salt, Username) VALUES
('hashed_pw1', 'salt1', 'john_doe'),
('hashed_pw2', 'salt2', 'jane_smith'),
('hashed_pw3', 'salt3', 'alice_nguyen'),
('hashed_pw4', 'salt4', 'bob_tran'),
('hashed_pw5', 'salt5', 'emma_phan');

-- Insert Users (Account link not in schema; assuming avatar will reference Metadata later)
INSERT INTO `User` (Email, Phone, Username, Name, avatar) VALUES
('john@example.com', '1234567890', 'john_doe', 'John Doe', NULL),
('jane@example.com', '2345678901', 'jane_smith', 'Jane Smith', NULL),
('alice@example.com', '3456789012', 'alice_nguyen', 'Alice Nguyen', NULL),
('bob@example.com', '4567890123', 'bob_tran', 'Bob Tran', NULL),
('emma@example.com', '5678901234', 'emma_phan', 'Emma Phan', NULL);

-- Insert Todos
INSERT INTO Todo (Description, Status, UpdateTime, CreateTime, Deadline, DeletedDate, Label, TagColor, workspace, TotalReact, Create_By, Assignee, Title)
VALUES
('Finish project report', 'Open', NOW(), NOW(), '2025-06-20', NULL, 'Work', 'Blue', 'ProjectA', 0, 1, 2, 'Project Report'),
('Buy groceries', 'Done', NOW(), NOW(), '2025-06-14', NULL, 'Personal', 'Green', 'Home', 3, 2, 2, 'Grocery List'),
('Plan team meeting', 'In Progress', NOW(), NOW(), '2025-06-16', NULL, 'Work', 'Red', 'ProjectB', 1, 3, 4, 'Team Meeting'),
('Write blog post', 'Open', NOW(), NOW(), '2025-06-18', NULL, 'Creative', 'Yellow', 'Blog', 2, 4, 3, 'Blog Draft'),
('Renew gym membership', 'Pending', NOW(), NOW(), '2025-06-15', NULL, 'Health', 'Purple', 'Fitness', 0, 5, 1, 'Gym Membership');

-- Insert Metadata (for Todo attachments or avatars)
INSERT INTO Metadata (FileName, FileExtension, FileSize, Height, Width, URL, Todo_id)
VALUES
('report.pdf', 'pdf', 102400, NULL, NULL, 'https://files.com/report.pdf', 1),
('groceries.jpg', 'jpg', 204800, 1080, 720, 'https://files.com/groceries.jpg', 2),
('meeting_notes.docx', 'docx', 51200, NULL, NULL, 'https://files.com/meeting.docx', 3),
('blog_draft.txt', 'txt', 1024, NULL, NULL, 'https://files.com/blog.txt', 4),
('gym_invoice.pdf', 'pdf', 20480, NULL, NULL, 'https://files.com/gym.pdf', 5);

-- Update avatar field in User table (linking to metadata)
UPDATE `User` SET avatar = 1 WHERE User_Id = 1;
UPDATE `User` SET avatar = 2 WHERE User_Id = 2;
UPDATE `User` SET avatar = 3 WHERE User_Id = 3;
UPDATE `User` SET avatar = 4 WHERE User_Id = 4;
UPDATE `User` SET avatar = 5 WHERE User_Id = 5;

-- Insert Comments
INSERT INTO Comment (Todo_Id, User_Id, Content) VALUES
(1, 2, 'I will review this tonight.'),
(2, 1, 'Nice shopping list.'),
(3, 4, 'Can we move it to Friday?'),
(4, 3, 'Looks great, maybe add some images.'),
(5, 5, 'Don’t forget to print the invoice.');

-- Insert Reactions
INSERT INTO Reaction (Todo_Id, User_Id, React) VALUES
(1, 3, 'Like'),
(2, 1, 'Love'),
(3, 5, 'Wow'),
(4, 2, 'Like'),
(5, 4, 'Dislike');



