ALTER TABLE todo_items
    ADD COLUMN IF NOT EXISTS due_date TIMESTAMP WITH TIME ZONE;