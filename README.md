# Гайд по проекту

Проект собран в виде упрощенного go-standard-project-layout.

Модели доменной области лежат в отдельном пакете models, бизнес-логика
лежит в каждой реализации datastore умышленно, так как нет смысла 
держать дополнительный абстрактный слой для данного задания.

Пакет datastore предоставляет единый интерфейс управления и взаимодействия
с данными, реализовано 2 имплементации данного интерфейса:

- Thread safe "in memory" OrderBook
- Clickhouse OrderBook
