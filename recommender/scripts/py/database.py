from pathlib import Path
from typing import List

import pandas as pd
import sqlalchemy
from pandas import DataFrame
from sqlalchemy import create_engine
from sqlalchemy.engine import Connection


def get_db_connection():
    return create_engine('postgresql+psycopg2://root:root@localhost/movie_match', pool_recycle=3600).connect()


def get_voted_media(conn: Connection, user_id: str) -> DataFrame:
    f = open(Path(__file__).parent / '../sql/media_voted.sql', 'r')

    return pd.read_sql_query(
        sqlalchemy.text(f.read()),
        conn,
        params={
            'user_id': user_id
        }
    )


def get_all_media(conn: Connection) -> DataFrame:
    f = open(Path(__file__).parent / '../sql/media_all.sql', 'r')

    return pd.read_sql_query(
        f.read(),
        conn
    )


def get_all_user_ids(conn: Connection) -> List[str]:
    return pd.read_sql_query(
        'select id from users',
        conn
    )['id'].values.tolist()
