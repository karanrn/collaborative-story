-- Sentence holds words (sentence_id is unique)
CREATE TABLE sentence (
    sentence_id int,
    word varchar(255),
    created_at timestamp
);

-- Paragraph holds size (start - end sentence) of the paragraph
CREATE TABLE paragraph (
    paragraph_id int PRIMARY KEY,
    start_sentence int,
    end_sentence int,
    created_at timestamp,
    updated_at timestamp
);

-- Story holds size of the story in terms of paragraph
CREATE TABLE story (
    story_id int PRIMARY KEY,
    title varchar(255),
    start_paragraph int,
    end_paragraph int,
    created_at timestamp,
    updated_at timestamp
);