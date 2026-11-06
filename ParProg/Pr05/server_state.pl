% Объявляем правила изменяемыми
:- dynamic cpu_load/2.
:- dynamic memory_load/2.
:- dynamic disk_free/2.
:- dynamic service_down/2.

% Факты
server(server1).
server(server2).
server(server3).
server(server4).
server(input_server).

cpu_load(server1, 92).
cpu_load(server2, 25).
cpu_load(server3, 55).
cpu_load(server4, 60).
cpu_load(input_server, 0).

memory_load(server1, 88).
memory_load(server2, 40).
memory_load(server3, 60).
memory_load(server4, 65).
memory_load(input_server, 0).

disk_free(server1, 40).
disk_free(server2, 70).
disk_free(server3, 5).
disk_free(server4, 35).
disk_free(input_server, 0).

service_down(server3, nginx).
service_down(server4, postgresql).

% Перезапись параметров
set_cpu(Server, Value) :-
    % Очищаем предыдущее значение
    retractall(cpu_load(Server, _)),
    % Добавляем новое значение
    assertz(cpu_load(Server, Value)).

set_memory(Server, Value) :-
    % Очищаем предыдущее значение
    retractall(memory_load(Server, _)),
    % Добавляем новое значение
    assertz(memory_load(Server, Value)).

set_disk_free(Server, Value) :-
    % Очищаем предыдущее значение
    retractall(disk_free(Server, _)),
    % Добавляем новое значение
    assertz(disk_free(Server, Value)).

% Простые правила
warn_cpu(Server) :- 
    cpu_load(Server, Percent),
    Percent >= 70,
    \+ critical_cpu(Server).

warn_memory(Server) :- 
    memory_load(Server, Percent),
    Percent >= 70,
    \+ critical_memory(Server).

warn_disk(Server) :- 
    disk_free(Server, Percent),
    Percent =< 20,
    \+ critical_disk(Server).

critical_cpu(Server) :- 
    cpu_load(Server, Percent),
    Percent >= 85.

critical_memory(Server) :- 
    memory_load(Server, Percent),
    Percent >= 85.

critical_disk(Server) :- 
    disk_free(Server, Percent),
    Percent =< 10.

has_down_service(Server) :- 
    service_down(Server, _).

overloaded(Server) :-
    critical_cpu(Server),
    critical_memory(Server).

resources_ok(Server) :-
    \+ warn_cpu(Server),
    \+ critical_cpu(Server),
    \+ warn_memory(Server),
    \+ critical_memory(Server),
    \+ warn_disk(Server),
    \+ critical_disk(Server).

% Итоговые правила
server_state(Server, critical) :-
    critical_disk(Server);
    overloaded(Server).

server_state(Server, warn) :-
    warn_disk(Server);
    warn_cpu(Server);
    warn_memory(Server);
    has_down_service(Server).

server_state(Server, normal) :-
    resources_ok(Server),
    \+ has_down_service(Server).

% Рекомендации
action(Server, 'Срочно добавить CPU/RAM или перераспределить нагрузку.') :-
    server_state(Server, critical),
    overloaded(Server).

action(Server, 'Срочно очистить диск или расширить хранилище.') :-
    server_state(Server, critical),
    critical_disk(Server).

action(Server, 'Поднять сервис и проверить ресурсы.') :-
    server_state(Server, warn),
    has_down_service(Server).

action(Server, 'Усилить мониторинг ресурсов.') :-
    server_state(Server, warn),
    \+ has_down_service(Server).

action(Server, 'Действий не требуется.') :-
    server_state(Server, normal).