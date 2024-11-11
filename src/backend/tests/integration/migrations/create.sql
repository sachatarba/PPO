-- Создание тестовых данных для таблицы gyms
INSERT INTO public.gyms (id, name, phone, city, addres, is_chain)
VALUES 
    ('a470d040-9495-418b-b2ab-dbd78bc9d0b3', 'FitLife Gym', '+7-999-99-99-99', 'New York', '123 Fitness St', TRUE),
    ('1fc09934-406d-4bdd-b8bf-ca99bafe8b30', 'PowerHouse', '+7-999-99-99-99', 'Los Angeles', '456 Power Rd', FALSE),
    ('08012002-95db-4a18-8e2b-6edb13d9d738', 'Gym Hero', '+7-999-99-99-99', 'Chicago', '789 Muscle Ave', TRUE);

-- Создание тестовых данных для таблицы trainers
INSERT INTO public.trainers (id, fullname, email, phone, qualification, unit_price)
VALUES 
    (gen_random_uuid(), 'John Doe', 'john.doe@example.com', '+7-999-99-99-99', 'Certified Trainer', 50.0),
    (gen_random_uuid(), 'Jane Smith', 'jane.smith@example.com', '+7-999-99-99-99', 'Strength Coach', 70.0),
    (gen_random_uuid(), 'Alice Johnson', 'alice.j@example.com', '+7-999-99-99-99', 'Yoga Instructor', 60.0);

-- Создание тестовых данных для таблицы membership_types
INSERT INTO public.membership_types (id, type, description, price, days_duration, gym_id)
VALUES 
    (gen_random_uuid(), 'Monthly', 'Monthly membership', 30.0, 30, (SELECT id FROM public.gyms LIMIT 1 OFFSET 0)),
    (gen_random_uuid(), 'Quarterly', 'Quarterly membership', 80.0, 90, (SELECT id FROM public.gyms LIMIT 1 OFFSET 1)),
    (gen_random_uuid(), 'Yearly', 'Yearly membership', 300.0, 365, (SELECT id FROM public.gyms LIMIT 1 OFFSET 2));

-- Создание тестовых данных для таблицы clients
INSERT INTO public.clients (id, login, password, fullname, email, phone, birthdate)
VALUES 
    (gen_random_uuid(), 'jdoe', 'password123', 'John Doe', 'jdoe@example.com', '+7-999-99-99-99', '1990-01-01'),
    (gen_random_uuid(), 'jsmith', 'password456', 'Jane Smith', 'jsmith@example.com', '+7-999-99-99-99', '1985-05-15'),
    (gen_random_uuid(), 'ajohnson', 'password789', 'Alice Johnson', 'ajohnson@example.com', '+7-999-99-99-99', '1992-03-22');

-- Создание тестовых данных для таблицы client_memberships
INSERT INTO public.client_memberships (id, start_date, end_date, membership_type_id, client_id)
VALUES 
    (gen_random_uuid(), '2024-01-01', '2024-12-31', (SELECT id FROM public.membership_types LIMIT 1 OFFSET 0), (SELECT id FROM public.clients LIMIT 1 OFFSET 0)),
    (gen_random_uuid(), '2024-01-01', '2024-04-01', (SELECT id FROM public.membership_types LIMIT 1 OFFSET 1), (SELECT id FROM public.clients LIMIT 1 OFFSET 1)),
    (gen_random_uuid(), '2024-01-01', '2024-01-31', (SELECT id FROM public.membership_types LIMIT 1 OFFSET 2), (SELECT id FROM public.clients LIMIT 1 OFFSET 2));

-- Создание тестовых данных для таблицы equipment
INSERT INTO public.equipment (id, name, description, gym_id)
VALUES 
    (gen_random_uuid(), 'Treadmill', 'Advanced treadmill for cardio', (SELECT id FROM public.gyms LIMIT 1 OFFSET 0)),
    (gen_random_uuid(), 'Dumbbells', 'Set of dumbbells from 1 to 50 kg', (SELECT id FROM public.gyms LIMIT 1 OFFSET 1)),
    (gen_random_uuid(), 'Bench Press', 'Bench press station', (SELECT id FROM public.gyms LIMIT 1 OFFSET 2));

INSERT INTO public.trainings (id, title, description, training_type, trainer_id)
VALUES 
    (gen_random_uuid(), 'Yoga Class', 'Yoga session for flexibility', 'aerobic', (SELECT id FROM public.trainers LIMIT 1 OFFSET 0)),
    (gen_random_uuid(), 'Strength Training', 'Building muscle strength', 'strength', (SELECT id FROM public.trainers LIMIT 1 OFFSET 1)),
    (gen_random_uuid(), 'Cardio Blast', 'High-intensity cardio workout', 'anaerobic', (SELECT id FROM public.trainers LIMIT 1 OFFSET 2));

INSERT INTO public.schedules (id, day_of_the_week, start_time, end_time, client_id, training_id)
VALUES 
    (gen_random_uuid(), '2024-11-13', '2024-11-13 09:00:00+00', '2024-11-13 10:00:00+00', (SELECT id FROM public.clients LIMIT 1 OFFSET 0), (SELECT id FROM public.trainings LIMIT 1 OFFSET 0)),
    (gen_random_uuid(), '2024-11-14', '2024-11-14 11:00:00+00', '2024-11-14 12:00:00+00', (SELECT id FROM public.clients LIMIT 1 OFFSET 1), (SELECT id FROM public.trainings LIMIT 1 OFFSET 1)),
    (gen_random_uuid(), '2024-11-15', '2024-11-15 13:00:00+00', '2024-11-15 14:00:00+00', (SELECT id FROM public.clients LIMIT 1 OFFSET 2), (SELECT id FROM public.trainings LIMIT 1 OFFSET 2));

INSERT INTO public.gym_trainers (trainer_id, gym_id)
VALUES 
    ((SELECT id FROM public.trainers LIMIT 1 OFFSET 0), (SELECT id FROM public.gyms LIMIT 1 OFFSET 0)),
    ((SELECT id FROM public.trainers LIMIT 1 OFFSET 1), (SELECT id FROM public.gyms LIMIT 1 OFFSET 1)),
    ((SELECT id FROM public.trainers LIMIT 1 OFFSET 2), (SELECT id FROM public.gyms LIMIT 1 OFFSET 2));

COMMIT;
