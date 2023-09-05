select m.id                                           as media_id,
       m.title                                        as media_title,
       m.summary                                      as media_summary,
       m.rating                                       as media_rating,
       m.release_date                                 as media_release_date,
       m.runtime                                      as media_runtime,
       m.type                                         as media_type,
       extract(epoch from m.release_date)             as media_release_date,
       max(case when mg.row_num = 1 then mg.name end) as media_genre_0,
       max(case when mg.row_num = 2 then mg.name end) as media_genre_1,
       max(case when mg.row_num = 3 then mg.name end) as media_genre_2,
       max(case when mg.row_num = 4 then mg.name end) as media_genre_3
from media m
         left join (select mg.media_id,
                      g.name,
                      row_number() over (partition by mg.media_id order by g.name) as row_num
               from media_genres mg
                        join genres g on g.id = mg.genre_id) mg on mg.media_id = m.id
group by m.id;
