import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import { TwoFactorForm } from '@/components/features/2fa/components/TwoFactorForm';

describe('LOGIN-02 - ingreso y validación visual de TOTP', () => {
  it('envía exactamente los seis dígitos ingresados', async () => {
    const onSubmit = vi.fn().mockResolvedValue(undefined);
    const user = userEvent.setup();
    render(<TwoFactorForm onSubmit={onSubmit} isSubmitting={false} />);

    const input = screen.getByRole('textbox');
    await user.type(input, '123456');
    await user.click(screen.getByRole('button', { name: 'Verificar código' }));

    expect(onSubmit).toHaveBeenCalledWith('123456');
  });

  it('acepta únicamente dígitos', async () => {
    const user = userEvent.setup();
    render(<TwoFactorForm onSubmit={vi.fn()} isSubmitting={false} />);

    const input = screen.getByRole('textbox');
    await user.type(input, 'ABCDEF');

    expect(input).toHaveValue('');
    expect(screen.getByRole('button', { name: 'Verificar código' })).toBeDisabled();
  });

  it('presenta errores y bloquea un nuevo envío durante la validación', () => {
    render(<TwoFactorForm onSubmit={vi.fn()} isSubmitting error="El código es incorrecto." />);

    expect(screen.getByRole('alert')).toHaveTextContent('El código es incorrecto.');
    expect(screen.getByRole('button', { name: 'Verificar código' })).toBeDisabled();
  });
});
