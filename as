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
      "evaluation": ,
      "score": ,
      "video_link": "https://example.com/video1",
      "emotion_results": [

      ]
    },
    {
      "question": "Describe a challenging project you have worked on.",
      "evaluation": ,
      "score": ,
      "video_link": "https://example.com/video2",
      "emotion_results": [

      ]
    }
  ],
  "score": 17,
  "video": "https://example.com/interview_video"
}