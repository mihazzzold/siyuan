// В self-hosted форке платные барьеры отключены: любая функция считается доступной,
// а пользователь — «оплаченным». Параметры сохранены ради совместимости сигнатур.

/* eslint-disable @typescript-eslint/no-unused-vars */
export const needSubscribe = (_tip: string = "") => {
    return false;
};

/**
 * 判断是否可以使用第三方同步：本 fork 中始终允许
 */
export const isPaidUser = () => {
    return true;
};
