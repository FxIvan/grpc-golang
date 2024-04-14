CREATE EXTENSION IF NO EXISTS "uuid-ossp";
 /*
 uuid_generate_v4() lo utilizamos gracias a uuidossp
 IF NO EXISTS es para que no de error si ya existe la extension
 */
CREATE TABLE tasks{
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    command TEXT NOT NULL,
    scheduled_at TIMESTAMP NOT NULL,
    picked_at TIMESTAMP,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    failed_at TIMESTAMP,
}
/*
UUID (Identificador único universal) es un número de 128 bits que se usa para identificar información en sistemas de computación.
NOT NULL vendria ah hacer como el required de mongoose
TIMETAMP es un tipo de dato que almacena la fecha y hora en la que se inserto el registro
*/

CREATE INDEX tasks_scheduled_at_idx ON tasks(scheduled_at);
/*
Creamos un indice llamada tasks_scheduled_at_idx en la tabla tasks con la columna scheduled_at
es decir que va van a ordenar los registros por la columna scheduled_at que es el horario se ejecucion
*/