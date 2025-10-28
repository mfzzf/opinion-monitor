# Testing Guide

## Manual Testing Checklist

### Backend Testing

#### 1. Setup Verification
- [ ] MySQL database created
- [ ] Backend config.yaml configured
- [ ] FFmpeg installed and accessible
- [ ] Server starts without errors

#### 2. Authentication Tests

```bash
# Register a new user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'

# Expected: 201 Created with token and user object

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'

# Expected: 200 OK with token and user object

# Get current user (replace TOKEN with actual token)
curl -X GET http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer TOKEN"

# Expected: 200 OK with user object
```

#### 3. Video Upload Tests

```bash
# Upload a video (replace TOKEN and VIDEO_PATH)
curl -X POST http://localhost:8080/api/videos/upload \
  -H "Authorization: Bearer TOKEN" \
  -F "videos=@/path/to/video.mp4"

# Expected: 201 Created with video objects

# List videos
curl -X GET http://localhost:8080/api/videos \
  -H "Authorization: Bearer TOKEN"

# Expected: 200 OK with videos array

# Get specific video (replace VIDEO_ID)
curl -X GET http://localhost:8080/api/videos/VIDEO_ID \
  -H "Authorization: Bearer TOKEN"

# Expected: 200 OK with video object
```

#### 4. Processing Pipeline Tests

- [ ] Video uploaded successfully
- [ ] Job created in database
- [ ] Video status changes to "processing"
- [ ] Cover frame extracted (check uploads directory)
- [ ] Text extracted from cover
- [ ] Sentiment analysis completed
- [ ] Report saved to database
- [ ] Video status changes to "completed"

Check backend logs for processing steps.

#### 5. Report Tests

```bash
# Get report by video ID
curl -X GET http://localhost:8080/api/reports/VIDEO_ID \
  -H "Authorization: Bearer TOKEN"

# Expected: 200 OK with report object containing:
# - cover_text
# - sentiment_score (0-1)
# - sentiment_label (positive/neutral/negative)
# - key_topics (JSON array)
# - risk_level (low/medium/high)
# - detailed_analysis
# - recommendations (JSON array)
```

### Frontend Testing

#### 1. Authentication Flow
- [ ] Visit http://localhost:3000
- [ ] Redirects to /login if not authenticated
- [ ] Can register new account
- [ ] Form validation works (email, password length, etc.)
- [ ] Can login with credentials
- [ ] Redirects to /videos after login
- [ ] Can logout
- [ ] Session persists on page refresh
- [ ] Redirects to /login after logout

#### 2. Video Upload
- [ ] Navigate to /upload
- [ ] Can select files via button
- [ ] Can drag and drop files
- [ ] Only video files accepted
- [ ] Can remove files before upload
- [ ] Upload progress shown
- [ ] Redirects to /videos after upload
- [ ] Uploaded videos appear in list

#### 3. Video List
- [ ] All videos displayed
- [ ] Status badges show correct state
- [ ] Can filter by status
- [ ] Pagination works (if > 12 videos)
- [ ] Auto-refresh updates status
- [ ] Can delete videos
- [ ] "View Report" button only for completed videos

#### 4. Report Detail
- [ ] Report page loads for completed video
- [ ] Video cover image displayed
- [ ] Extracted text shown
- [ ] Sentiment score displayed with progress bar
- [ ] Risk level badge shown
- [ ] Key topics displayed as tags
- [ ] Detailed analysis text shown
- [ ] Recommendations listed
- [ ] All metadata correct
- [ ] Back button returns to videos

#### 5. Error Handling
- [ ] Invalid login shows error
- [ ] Upload without file shows error
- [ ] Network errors shown gracefully
- [ ] 401 errors redirect to login
- [ ] Report page handles missing reports

### Integration Tests

#### End-to-End Test Scenario

1. **Setup**
   - Start backend
   - Start frontend
   - Prepare test video file

2. **User Registration**
   - Register new user
   - Verify email validation
   - Verify password requirements
   - Confirm redirect to videos

3. **Video Upload**
   - Upload single video
   - Verify upload progress
   - Confirm redirect to videos
   - Check video appears with "pending" status

4. **Processing**
   - Wait for status to change to "processing"
   - Monitor backend logs for:
     - Job picked up
     - Cover extraction
     - Text extraction API call
     - Sentiment analysis API call
     - Report saved
   - Wait for status to change to "completed"

5. **Report Viewing**
   - Click "View Report"
   - Verify all report fields populated
   - Check sentiment score is 0-1
   - Verify sentiment label matches score
   - Confirm risk level is low/medium/high
   - Check key topics are relevant
   - Verify recommendations are actionable

6. **Batch Upload**
   - Upload 3-5 videos at once
   - Verify all appear in list
   - Confirm parallel processing
   - Check all complete successfully

7. **Cleanup**
   - Delete uploaded videos
   - Verify files removed from disk
   - Logout

### Performance Tests

#### Load Testing
- Upload 10 videos simultaneously
- Check worker pool handles concurrency
- Verify no deadlocks or crashes
- Confirm all videos process successfully

#### Stress Testing
- Upload large video (near 500MB limit)
- Verify upload completes
- Check processing handles large files
- Monitor memory usage

### Database Tests

```sql
-- Check user created
SELECT * FROM users WHERE email = 'test@example.com';

-- Check video record
SELECT * FROM videos WHERE user_id = YOUR_USER_ID;

-- Check job created
SELECT * FROM jobs WHERE video_id = YOUR_VIDEO_ID;

-- Check report generated
SELECT * FROM reports WHERE video_id = YOUR_VIDEO_ID;

-- Verify foreign key relationships
SELECT v.id, v.original_filename, v.status, j.status as job_status, r.sentiment_label
FROM videos v
LEFT JOIN jobs j ON v.id = j.video_id
LEFT JOIN reports r ON v.id = r.video_id
WHERE v.user_id = YOUR_USER_ID;
```

### API Tests

Use the provided curl commands or tools like:
- Postman
- Insomnia
- HTTPie

### Common Issues to Test

1. **Authentication**
   - Expired tokens
   - Invalid tokens
   - Missing tokens

2. **File Upload**
   - Non-video files
   - Files too large
   - Empty files
   - Corrupt files

3. **Processing**
   - Videos without text in cover
   - Very long videos
   - Unusual formats

4. **Edge Cases**
   - Empty video list
   - Failed processing
   - Network timeouts
   - API rate limits

## Automated Testing

For production, consider adding:

### Backend Tests
```go
// Example test structure
func TestUserRegistration(t *testing.T) {
    // Test user registration endpoint
}

func TestVideoUpload(t *testing.T) {
    // Test video upload endpoint
}

func TestVideoProcessing(t *testing.T) {
    // Test full processing pipeline
}
```

### Frontend Tests
```typescript
// Example with Jest/React Testing Library
describe('Login Page', () => {
  it('should render login form', () => {
    // Test login form rendering
  });

  it('should validate email format', () => {
    // Test email validation
  });
});
```

## Test Data

Sample test cases:
- **Short video (< 10s)**: Fast processing test
- **Medium video (30s-1m)**: Normal use case
- **Video with text overlay**: Text extraction test
- **Video without text**: Edge case handling
- **Multiple videos**: Batch processing test

## Success Criteria

The system passes testing if:
- ✅ All API endpoints respond correctly
- ✅ Authentication works end-to-end
- ✅ Videos upload successfully
- ✅ Processing completes without errors
- ✅ Reports are generated accurately
- ✅ Frontend displays all data correctly
- ✅ Error handling is graceful
- ✅ Performance is acceptable (< 30s for typical video)

## Notes

- First API call may be slower (cold start)
- OpenAI API rate limits may affect testing
- Large batch uploads require sufficient API credits
- Processing time depends on video size and API response time

