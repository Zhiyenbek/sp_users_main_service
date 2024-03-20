SELECT p.recruiter_public_id, COUNT(*) as interview_count
FROM interviews i
INNER JOIN user_interviews ui ON ui.interview_id = i.id
INNER JOIN positions p ON p.id = ui.position_id
INNER JOIN recruiters r ON p.recruiter_public_id = r.public_id
WHERE p.recruiter_public_id = $1
GROUP BY p.recruiter_public_id;

{
  "questions": [
    {
      "question": "What is your experience with object-oriented programming?",
      "evaluation": "Good",
      "score": 8,
      "video_link": "https://example.com/video1",
      "emotion_results": [
        {
          "emotion": "Happiness",
          "exact_time": 24.5,
          "duration": 10.2
        },
        {
          "emotion": "Neutral",
          "exact_time": 36.2,
          "duration": 5.7
        }
      ]
    },
    {
      "question": "Describe a challenging project you have worked on.",
      "evaluation": "Excellent performance with exceptional problem-solving skills",
      "score": 9,
      "video_link": "https://example.com/video2",
      "emotion_results": [
        {
          "emotion": "Confidence",
          "exact_time": 45.8,
          "duration": 8.5
        },
        {
          "emotion": "Determination",
          "exact_time": 56.3,
          "duration": 7.1
        }
      ]
    }
  ],
  "score": 17,
  "video": "https://example.com/interview_video"
}