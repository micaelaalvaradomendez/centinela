-- Seed idempotente para demo local de El Centinela.
-- Crea/actualiza usuarios ADMIN y OPERATOR con contraseña: Prueba123$
-- y asigna permisos sobre instancias del simulador de Proxmox VE.

DO $$
DECLARE
    v_org_id uuid;
    v_admin_id uuid;
    v_op1_id uuid;
    v_op2_id uuid;
    v_hash text := crypt('Prueba123$', gen_salt('bf', 12));
BEGIN
    -- 1. Reutilizar organización existente o crear una por defecto
    SELECT id INTO v_org_id FROM organizaciones ORDER BY nombre_org LIMIT 1;
    IF v_org_id IS NULL THEN
        v_org_id := uuid_generate_v7();
        INSERT INTO organizaciones (id, nombre_org, creado_por)
        VALUES (v_org_id, 'El Centinela Demo', '00000000-0000-0000-0000-000000000000');
    END IF;

    -- 2. Usuario ADMIN: admin@elcentinela.com / Prueba123$
    SELECT id INTO v_admin_id
    FROM usuarios
    WHERE lower(btrim(email_usuario)) = 'admin@elcentinela.com'
      AND eliminado_en IS NULL
    LIMIT 1;

    IF v_admin_id IS NULL THEN
        v_admin_id := uuid_generate_v7();
        INSERT INTO usuarios (
            id, organizacion_id, nombre_completo, nombre_usuario, email_usuario,
            contrasena_hash, rol, activo, cambio_contrasena, totp_vinculado, fecha_creacion
        ) VALUES (
            v_admin_id, v_org_id, 'Administrador Demo', 'admin', 'admin@elcentinela.com',
            v_hash, 'ADMIN', true, false, false, now()
        );
    ELSE
        UPDATE usuarios
        SET contrasena_hash = v_hash,
            rol = 'ADMIN',
            activo = true,
            cambio_contrasena = false
        WHERE id = v_admin_id;
    END IF;

    -- Asegurar creado_por de la organización apuntando al admin
    UPDATE organizaciones
    SET creado_por = v_admin_id
    WHERE id = v_org_id
      AND creado_por = '00000000-0000-0000-0000-000000000000';

    -- 3. Usuario OPERATOR 1: operador@elcentinela.com / Prueba123$
    SELECT id INTO v_op1_id
    FROM usuarios
    WHERE lower(btrim(email_usuario)) = 'operador@elcentinela.com'
      AND eliminado_en IS NULL
    LIMIT 1;

    IF v_op1_id IS NULL THEN
        v_op1_id := uuid_generate_v7();
        INSERT INTO usuarios (
            id, organizacion_id, nombre_completo, nombre_usuario, email_usuario,
            contrasena_hash, rol, activo, cambio_contrasena, totp_vinculado, fecha_creacion
        ) VALUES (
            v_op1_id, v_org_id, 'Operador Principal', 'operador', 'operador@elcentinela.com',
            v_hash, 'OPERATOR', true, false, false, now()
        );
    ELSE
        UPDATE usuarios
        SET contrasena_hash = v_hash,
            rol = 'OPERATOR',
            activo = true,
            cambio_contrasena = false
        WHERE id = v_op1_id;
    END IF;

    -- 4. Usuario OPERATOR 2: operador2@elcentinela.com / Prueba123$
    SELECT id INTO v_op2_id
    FROM usuarios
    WHERE lower(btrim(email_usuario)) = 'operador2@elcentinela.com'
      AND eliminado_en IS NULL
    LIMIT 1;

    IF v_op2_id IS NULL THEN
        v_op2_id := uuid_generate_v7();
        INSERT INTO usuarios (
            id, organizacion_id, nombre_completo, nombre_usuario, email_usuario,
            contrasena_hash, rol, activo, cambio_contrasena, totp_vinculado, fecha_creacion
        ) VALUES (
            v_op2_id, v_org_id, 'Operador Secundario', 'operador2', 'operador2@elcentinela.com',
            v_hash, 'OPERATOR', true, false, false, now()
        );
    ELSE
        UPDATE usuarios
        SET contrasena_hash = v_hash,
            rol = 'OPERATOR',
            activo = true,
            cambio_contrasena = false
        WHERE id = v_op2_id;
    END IF;

    -- 5. Permisos sobre instancias del simulador Proxmox para los operadores
    INSERT INTO permisos_instancia (id, usuario_id, vmid_proxmox, nivel_acceso)
    VALUES
        (uuid_generate_v7(), v_op1_id, 110, 'FULL_ACCESS'),
        (uuid_generate_v7(), v_op1_id, 201, 'FULL_ACCESS'),
        (uuid_generate_v7(), v_op1_id, 9000, 'FULL_ACCESS'),
        (uuid_generate_v7(), v_op1_id, 9001, 'READ_ONLY'),
        (uuid_generate_v7(), v_op1_id, 9002, 'READ_ONLY'),
        (uuid_generate_v7(), v_op1_id, 9003, 'FULL_ACCESS'),
        (uuid_generate_v7(), v_op2_id, 110, 'READ_ONLY'),
        (uuid_generate_v7(), v_op2_id, 9002, 'FULL_ACCESS')
    ON CONFLICT (usuario_id, vmid_proxmox) DO NOTHING;
END $$;
