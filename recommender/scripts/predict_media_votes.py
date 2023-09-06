from sqlalchemy import text

from database import get_db_connection, get_all_user_ids
from predict import predict_votes

conn = get_db_connection()

user_ids = get_all_user_ids(conn)

for user_id in get_all_user_ids(conn):
    p_media = predict_votes(conn, user_id)

    for idx, row in p_media.iterrows():
        print(f'inserting vote prediction {row["vote_type"]} for {row["media_id"]}.')

        conn.execute(
            text(
                'insert into media_user_vote_prediction ('
                'media_id, user_id, predicted_vote, created_at, updated_at'
                ') values ('
                ':media_id, :user_id, :predicted_vote, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP'
                ')'
                ' on conflict (media_id, user_id) do update set'
                ' predicted_vote = :predicted_vote, updated_at = CURRENT_TIMESTAMP'
            ),
            {
                'media_id': row['media_id'],
                'user_id': user_id,
                'predicted_vote': row['vote_type'],
            }
        )

        conn.commit()

conn.close()
