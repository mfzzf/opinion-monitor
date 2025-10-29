#!/bin/bash
# Script to verify database migration for audio fields

echo "========================================"
echo "Database Migration Verification"
echo "========================================"
echo ""

# Database credentials from config.yaml
DB_HOST="localhost"
DB_PORT="3306"
DB_NAME="opinion_monitor"
DB_USER="root"
DB_PASS="Mei123123!"

echo "Checking videos table structure..."
mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "DESCRIBE videos;" 2>/dev/null | grep -E "(audio_path|transcript_text)"

if [ $? -eq 0 ]; then
    echo "✓ Videos table has new fields"
else
    echo "✗ Videos table missing new fields"
    echo "  Run: go run cmd/server/main.go (GORM will auto-migrate)"
fi

echo ""
echo "Checking reports table structure..."
mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "DESCRIBE reports;" 2>/dev/null | grep "transcript_text"

if [ $? -eq 0 ]; then
    echo "✓ Reports table has new field"
else
    echo "✗ Reports table missing new field"
    echo "  Run: go run cmd/server/main.go (GORM will auto-migrate)"
fi

echo ""
echo "========================================"
echo "Note: GORM AutoMigrate will add these fields"
echo "automatically when you start the backend server."
echo ""
echo "Start the backend with:"
echo "  cd backend && go run cmd/server/main.go"
echo "========================================"

