CREATE TABLE users
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    user_name   VARCHAR(255) UNIQUE NOT NULL,
    password    VARCHAR(255)        NOT NULL,
    user_type   ENUM ('ADMIN', 'STUDENT') DEFAULT 'STUDENT',
    create_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    delete_date TIMESTAMP NULL
);

CREATE TABLE schools
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    region      VARCHAR(255) NOT NULL,
    user_id     INT          NOT NULL,
    create_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    delete_date TIMESTAMP NULL,
    UNIQUE KEY `unique_school_region` (`name`, `region`),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE news
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    user_id     INT          NOT NULL,
    school_id   INT          NOT NULL,
    create_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    delete_date TIMESTAMP NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (school_id) REFERENCES schools (id)
);

CREATE TABLE subscriptions
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    user_id     INT          NOT NULL,
    school_id   INT          NOT NULL,
    create_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    delete_date TIMESTAMP NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (school_id) REFERENCES schools (id)
);

-- 유저 생성 (관리자 3, 학생 1)
INSERT INTO users (user_name, password, user_type) VALUES ('classting_admin_1', '$2a$10$Q2w/PDGj.iS1vKCMRnIJ1uuTG7jzQn13Mu4IplYgwmyTtod7yCHAu', 'ADMIN');
INSERT INTO users (user_name, password, user_type) VALUES ('classting_admin_2', '$2a$10$Q2w/PDGj.iS1vKCMRnIJ1uuTG7jzQn13Mu4IplYgwmyTtod7yCHAu', 'ADMIN');
INSERT INTO users (user_name, password, user_type) VALUES ('classting_admin_3', '$2a$10$Q2w/PDGj.iS1vKCMRnIJ1uuTG7jzQn13Mu4IplYgwmyTtod7yCHAu', 'ADMIN');
INSERT INTO users (user_name, password, user_type) VALUES ('classting_student_1', '$2a$10$kxG7/4JAyyRkGcn7ZkgZFu2RpUpRtf/CnA.L4p/L7JaAbth/iQKUG', 'STUDENT');
-- 비어있는 계정
INSERT INTO users (user_name, password, user_type) VALUES ('empty_classting_admin', '$2a$10$Q2w/PDGj.iS1vKCMRnIJ1uuTG7jzQn13Mu4IplYgwmyTtod7yCHAu', 'ADMIN');
INSERT INTO users (user_name, password, user_type) VALUES ('empty_classting_student', '$2a$10$Q2w/PDGj.iS1vKCMRnIJ1uuTG7jzQn13Mu4IplYgwmyTtod7yCHAu', 'STUDENT');
-- 학교 생성 (관리자 1, 2, 3)
INSERT INTO schools (name, region, user_id) VALUES ('admin_1_뉴스가_있는_클래스팅 서울학교', '서울', 1);
INSERT INTO schools (name, region, user_id) VALUES ('admin_1_뉴스가_있는_클래스팅 대전학교', '대전', 1);
INSERT INTO schools (name, region, user_id) VALUES ('admin_2_뉴스가_있는_클래스팅 다른 인천학교', '인천', 2);
-- 페이지네이션용 데이터
INSERT INTO schools (name, region, user_id) VALUES ('admin_3_페이지네이션_확인_학교_1', '울산', 3);
INSERT INTO schools (name, region, user_id) VALUES ('admin_3_페이지네이션_확인_학교_2', '울산', 3);
INSERT INTO schools (name, region, user_id) VALUES ('admin_3_페이지네이션_확인_학교_3', '울산', 3);
INSERT INTO schools (name, region, user_id) VALUES ('admin_3_페이지네이션_확인_학교_4', '울산', 3);
INSERT INTO schools (name, region, user_id) VALUES ('admin_3_페이지네이션_확인_학교_5', '울산', 3);
INSERT INTO schools (name, region, user_id) VALUES ('admin_3_페이지네이션_확인_학교_6', '울산', 3);
INSERT INTO schools (name, region, user_id) VALUES ('admin_3_페이지네이션_확인_학교_7', '울산', 3);
INSERT INTO schools (name, region, user_id) VALUES ('admin_3_페이지네이션_확인_학교_8', '울산', 3);
INSERT INTO schools (name, region, user_id) VALUES ('admin_3_페이지네이션_확인_학교_9', '울산', 3);
INSERT INTO schools (name, region, user_id) VALUES ('admin_3_페이지네이션_확인_학교_10', '울산', 3);
INSERT INTO schools (name, region, user_id) VALUES ('admin_3_페이지네이션_확인_학교_11', '울산', 3);
INSERT INTO schools (name, region, user_id) VALUES ('admin_3_페이지네이션_확인_학교_12', '울산', 3);
INSERT INTO schools (name, region, user_id) VALUES ('admin_3_페이지네이션_확인_학교_13', '울산', 3);
INSERT INTO schools (name, region, user_id) VALUES ('admin_3_페이지네이션_확인_학교_14', '울산', 3);
INSERT INTO schools (name, region, user_id) VALUES ('admin_3_페이지네이션_확인_학교_15', '울산', 3);
INSERT INTO schools (name, region, user_id) VALUES ('admin_3_페이지네이션_확인_학교_16', '울산', 3);
INSERT INTO schools (name, region, user_id) VALUES ('admin_3_페이지네이션_확인_학교_17', '울산', 3);
INSERT INTO schools (name, region, user_id) VALUES ('admin_3_페이지네이션_확인_학교_18', '울산', 3);
INSERT INTO schools (name, region, user_id) VALUES ('admin_3_페이지네이션_확인_학교_19', '울산', 3);
INSERT INTO schools (name, region, user_id) VALUES ('admin_3_페이지네이션_확인_학교_20', '울산', 3);
-- 학교 구독 (학생 1)
INSERT INTO subscriptions (user_id, school_id) VALUES (4, 1);
INSERT INTO subscriptions (user_id, school_id) VALUES (4, 2);
INSERT INTO subscriptions (user_id, school_id) VALUES (4, 3);
INSERT INTO subscriptions (user_id, school_id) VALUES (4, 4);
INSERT INTO subscriptions (user_id, school_id) VALUES (4, 5);
INSERT INTO subscriptions (user_id, school_id) VALUES (4, 6);
INSERT INTO subscriptions (user_id, school_id) VALUES (4, 7);
INSERT INTO subscriptions (user_id, school_id) VALUES (4, 8);
INSERT INTO subscriptions (user_id, school_id) VALUES (4, 9);
INSERT INTO subscriptions (user_id, school_id) VALUES (4, 10);
INSERT INTO subscriptions (user_id, school_id) VALUES (4, 11);
INSERT INTO subscriptions (user_id, school_id) VALUES (4, 12);
INSERT INTO subscriptions (user_id, school_id) VALUES (4, 13);
INSERT INTO subscriptions (user_id, school_id) VALUES (4, 14);
INSERT INTO subscriptions (user_id, school_id) VALUES (4, 15);
INSERT INTO subscriptions (user_id, school_id) VALUES (4, 16);
INSERT INTO subscriptions (user_id, school_id) VALUES (4, 17);
INSERT INTO subscriptions (user_id, school_id) VALUES (4, 18);
INSERT INTO subscriptions (user_id, school_id) VALUES (4, 19);
INSERT INTO subscriptions (user_id, school_id) VALUES (4, 20);
INSERT INTO subscriptions (user_id, school_id) VALUES (4, 21);
INSERT INTO subscriptions (user_id, school_id) VALUES (4, 22);
INSERT INTO subscriptions (user_id, school_id) VALUES (4, 23);
-- 뉴스 생성
INSERT INTO news (title, user_id, school_id) VALUES ('admin_1_뉴스가_있는_클래스팅 서울학교_뉴스_1 1page', 1, 1);
INSERT INTO news (title, user_id, school_id) VALUES ('admin_1_뉴스가_있는_클래스팅 서울학교_뉴스_2 1page', 1, 1);
INSERT INTO news (title, user_id, school_id) VALUES ('admin_1_뉴스가_있는_클래스팅 서울학교_뉴스_3 1page', 1, 1);
INSERT INTO news (title, user_id, school_id) VALUES ('admin_1_뉴스가_있는_클래스팅 서울학교_뉴스_4 1page', 1, 1);
INSERT INTO news (title, user_id, school_id) VALUES ('admin_1_뉴스가_있는_클래스팅 서울학교_뉴스_5 1page', 1, 1);
INSERT INTO news (title, user_id, school_id) VALUES ('admin_1_뉴스가_있는_클래스팅 서울학교_뉴스_6 1page', 1, 1);
INSERT INTO news (title, user_id, school_id) VALUES ('admin_1_뉴스가_있는_클래스팅 서울학교_뉴스_7 1page', 1, 1);
INSERT INTO news (title, user_id, school_id) VALUES ('admin_1_뉴스가_있는_클래스팅 서울학교_뉴스_8 1page', 1, 1);
INSERT INTO news (title, user_id, school_id) VALUES ('admin_1_뉴스가_있는_클래스팅 서울학교_뉴스_9 1page', 1, 1);
INSERT INTO news (title, user_id, school_id) VALUES ('admin_1_뉴스가_있는_클래스팅 서울학교_뉴스_10 1page', 1, 1);
INSERT INTO news (title, user_id, school_id) VALUES ('admin_1_뉴스가_있는_클래스팅 서울학교_뉴스_11 2page', 1, 1);

INSERT INTO news (title, user_id, school_id) VALUES ('admin_1_뉴스가_있는_클래스팅 대전학교_뉴스_1', 1, 2);
INSERT INTO news (title, user_id, school_id) VALUES ('admin_1_뉴스가_있는_클래스팅 대전학교_뉴스_2', 1, 2);
INSERT INTO news (title, user_id, school_id) VALUES ('admin_1_뉴스가_있는_클래스팅 대전학교_뉴스_3', 1, 2);
INSERT INTO news (title, user_id, school_id) VALUES ('admin_1_뉴스가_있는_클래스팅 대전학교_뉴스_4', 1, 2);
INSERT INTO news (title, user_id, school_id) VALUES ('admin_1_뉴스가_있는_클래스팅 대전학교_뉴스_5', 1, 2);

INSERT INTO news (title, user_id, school_id) VALUES ('admin_2_뉴스가_있는_클래스팅 인천학교_뉴스_1', 2, 3);




