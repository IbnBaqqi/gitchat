-- +goose Up
CREATE TABLE pinned_conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    pinned_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, conversation_id)
);

CREATE INDEX idx_pinned_conversations_user ON pinned_conversations(user_id);
CREATE INDEX idx_pinned_conversations_conversation ON pinned_conversations(conversation_id);

-- +goose Down
DROP TABLE IF EXISTS pinned_conversations;