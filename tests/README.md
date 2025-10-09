# Tests Directory

This directory contains test files and documentation for the Flux Server WhatsApp Bot.

## Test Files

### `test_response_format.go`
Tests the AI response parsing functionality to ensure that Gemini API responses are correctly parsed into structured advice sections.

**Usage:**
```bash
go run tests/test_response_format.go
```

**What it tests:**
- Parsing of AI responses into structured advice sections
- Extraction of planting, irrigation, harvest, market, and general advice
- Handling of markdown formatting in AI responses
- Fallback behavior for malformed responses

## Documentation Files

### `test_multiple_crops.md`
Comprehensive documentation of the multiple crop registration feature implementation, including:
- Database schema changes
- Bot flow modifications
- API integration details
- Test scenarios

### `demo_flow.md`
Complete demo flow documentation showing the end-to-end user experience:
- Registration process with multiple crops
- AI advice generation
- Market insights
- Web app access
- All available commands

## Running Tests

To run any test file:
```bash
cd /path/to/flux-server
go run tests/[test_file_name].go
```

## Test Environment

Make sure you have the following environment variables set:
- `API_KEY`: Your Gemini API key for AI integration tests

## Notes

- Test files are standalone and don't require the full application to be running
- They use mock data where appropriate to avoid external dependencies
- All tests are designed to be run independently
