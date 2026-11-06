connection(d1,d2). 
connection(d2,d3).
connection(d3,d5).
connection(d1,d4).
connection(d4,d5).

% Прямая связь
path(A, B) :-
    connection(A, B).

% Косвеная связь
path(A, B) :-
    connection(A, Mid),
    path(Mid, B).