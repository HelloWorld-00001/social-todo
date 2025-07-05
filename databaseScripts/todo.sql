

CREATE TABLE `Account` (
    Id INT AUTO_INCREMENT PRIMARY KEY,
    `Password` VARCHAR(255),
    Salt VARCHAR(255),
    Username VARCHAR(100) UNIQUE
);

CREATE TABLE `User` (
    Id INT AUTO_INCREMENT PRIMARY KEY,
    Email VARCHAR(200) UNIQUE,
    Phone VARCHAR(20) UNIQUE,
    Username VARCHAR(100) UNIQUE,
    `Name` VARCHAR(100),
    avatar Int
);




CREATE TABLE Todo (
    Id INT AUTO_INCREMENT PRIMARY KEY,
    Description TEXT CHARACTER SET utf8mb4,
    `Status` VARCHAR(50),
    UpdateTime DATETIME,
    CreateTime DATETIME,
    Deadline DATETIME,
    DeletedDate DATETIME,      
    Label VARCHAR(100),
    TagColor VARCHAR(20),
    workspace VARCHAR(100),
    TotalReact INT,
    Create_By INT,
    Assignee INT,
    Title VARCHAR(255),
    CONSTRAINT fk_todo_creator FOREIGN KEY (Create_By) REFERENCES `User`(User_Id),
    CONSTRAINT fk_todo_assignee FOREIGN KEY (Assignee) REFERENCES `User`(User_Id)
);

CREATE TABLE Metadata (
    Id INT AUTO_INCREMENT PRIMARY KEY,
    FileName VARCHAR(255),
    FileExtension VARCHAR(20),
    FileSize INT,
    Height INT,
    Width INT,
    URL VARCHAR(500),
    Todo_id INT,
    CONSTRAINT fk_metadata_owner FOREIGN KEY (Todo_id) REFERENCES Todo(Todo_id)
);


CREATE TABLE Comment (
    Todo_Id INT,
    User_Id INT,
    Content TEXT CHARACTER SET utf8mb4,
    PRIMARY KEY (Todo_Id, User_Id),
    CONSTRAINT fk_comment_todo FOREIGN KEY (Todo_Id) REFERENCES ToDo(Todo_Id),
    CONSTRAINT fk_comment_user FOREIGN KEY (User_Id) REFERENCES `User`(User_Id)
);

CREATE TABLE Reaction (
    Todo_Id INT,
    User_Id INT,
    React ENUM('Like', 'Dislike', 'Love', 'Angry', 'Wow'),
    PRIMARY KEY (Todo_Id, User_Id),
    CONSTRAINT fk_reaction_todo FOREIGN KEY (Todo_Id) REFERENCES Todo(Todo_Id),
    CONSTRAINT fk_reaction_user FOREIGN KEY (User_Id) REFERENCES `User`(User_Id)
);


-- Link between User and Account
ALTER TABLE `User`
ADD CONSTRAINT fk_user_metadata FOREIGN KEY (avatar) REFERENCES Metadata(metadata_id);