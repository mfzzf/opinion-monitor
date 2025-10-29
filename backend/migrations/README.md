# Database Migrations

## Manual Migration

If you need to manually apply migrations, run:

```bash
mysql -u root -p opinion_monitor < migrations/add_audio_fields.sql
```

## Automatic Migration

The application uses GORM's AutoMigrate feature, which will automatically create the new columns when the server starts. You can simply restart the backend server:

```bash
cd backend
go run cmd/server/main.go
```

GORM will detect the new fields in the models and add them to the database automatically.

## Verification

To verify the migration was successful:

```sql
-- Check videos table structure
DESCRIBE videos;

-- Check reports table structure
DESCRIBE reports;
```

You should see:
- `videos.audio_path` (VARCHAR(500))
- `videos.transcript_text` (TEXT)
- `reports.transcript_text` (TEXT)

