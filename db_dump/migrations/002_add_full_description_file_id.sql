ALTER TABLE public.model
    ADD COLUMN IF NOT EXISTS full_description_file_id integer;

ALTER TABLE public.model
    DROP CONSTRAINT IF EXISTS model_full_description_file_id_fk;

ALTER TABLE public.model
    ADD CONSTRAINT model_full_description_file_id_fk
    FOREIGN KEY (full_description_file_id) REFERENCES public.file(id) ON DELETE SET NULL;
