Мы проверили формат, сейчас все должно считываться правильно, но есть 
нюанс относительно подсчета метрик для текстовой части.

Я загрузил наш файл в тетрадь Гугл-коллаба,
которую организаторы скидывали участникам для проверки. 
А так, как эталонного файла, относительно которого будет проходить проверка, 
у нас нет, его же загрузил в качестве эталона. Поэтому по идее метрики для 
классификационной и текстовой части должны давать наилучшее значение.
Для классификационной так и получилось, а вот для текстовой, получилось
nan на выходе.
Я изучил, почему так получается, подсчет для текстовой части происходит 
вот здесь:

```
def score_texts(pred_df, real_df):
    text_columns = ["address", "cam_id"]
    Y_hat_t = pred_df[text_columns].fillna("")
    print(Y_hat_t.head())
    print(Y_hat_t[text_columns])
    Y_t = real_df[text_columns].fillna("")
    Y_hat_t["cam_id"] = Y_hat_t["cam_id"].apply(str)
    Y_t["cam_id"] = Y_t["cam_id"].apply(str)
    ad_lens = Y_hat_t["address"].apply(len).values
    cid_lens = Y_hat_t["cam_id"].apply(len).values
    largest_lengths = np.array(list(zip(max(a,b) for a,b in zip(ad_lens, cid_lens))))
    print(largest_lengths)
    ad_levs = np.array(
        list(zip(distance(a, b) for a,b in zip(Y_hat_t["address"], Y_t["address"])))
    )
    cid_levs = np.array(
        list(zip(distance(a, b) for a,b in zip(Y_hat_t["cam_id"], Y_t["cam_id"])))
    )
    print(ad_levs)
    print(cid_levs)
    ad_levs = ad_levs/largest_lengths
    cid_levs = cid_levs/largest_lengths
    text_score = np.mean(np.array([1/(min(a,b)+10**-5) for a,b in zip(ad_levs, cid_levs)]))
    return(text_score)
```

Выяснилось, что там, где пропуски для cam_id и address, код заменяет их на пустые строки, 
потом для каждой строки вычисляется largest_lengths, и на это значение 
делятся получившиеся расстояния Левенштейна 
(`ad_levs = ad_levs/largest_lengths` и `cid_levs = cid_levs/largest_lengths`) для 
каждой строки, потом эти все значения усредняются. Но для пустых строк largest_lengths=0, 
поэтому получается, что для некоторых строк у вас происходит деление на ноль.
Поэтому нужно отдельно придумать, как работать с текстовой частью. 
Например, не учитывать при подсчете метрики