-- Migration to add audio_path, transcript_text fields to videos and reports tables

-- Add fields to videos table
ALTER TABLE videos 
ADD COLUMN IF NOT EXISTS audio_path VARCHAR(500) AFTER cover_path,
ADD COLUMN IF NOT EXISTS transcript_text TEXT AFTER audio_path;

-- Add field to reports table
ALTER TABLE reports 
ADD COLUMN IF NOT EXISTS transcript_text TEXT AFTER cover_text;

-- Add indexes for better query performance (optional)
CREATE INDEX IF NOT EXISTS idx_videos_audio_path ON videos(audio_path(255));

