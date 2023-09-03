from scripts.py.database import get_db_connection, get_all_user_ids
from scripts.py.predict import predict_votes

conn = get_db_connection()

user_ids = get_all_user_ids(conn)

for user_id in get_all_user_ids(conn):
    p_media = predict_votes(conn, user_id)

    for idx, row in p_media.iterrows():
        conn.execute(
            'insert into media_user_vote_prediction ('
            'media_id, user_id, predicted_vote, created_at, updated_at'
            ') values ('
            '%(media_id)s, %(user_id)s, %(predicted_vote)s, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP'
            ')'
            ' on conflict (media_id, user_id) do update set'
            ' predicted_vote = %(predicted_vote)s, updated_at = CURRENT_TIMESTAMP',
            {
                'media_id': row['media_id'],
                'user_id': user_id,
                'predicted_vote': row['vote_type'],
            }
        )
