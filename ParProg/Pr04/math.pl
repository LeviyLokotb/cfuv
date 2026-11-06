average3(A,B,C,M) :- M is (A+B+C)/3.

% Пустой список
count_negative([], 0).

% Вершина < 0
count_negative([Head|Tail], Count) :- 
    Head < 0,
    count_negative(Tail, CountOther),
    Count is 1 + CountOther.

% Вершина >= 0
count_negative([Head|Tail], Count) :- 
    Head >= 0,
    count_negative(Tail, Count).