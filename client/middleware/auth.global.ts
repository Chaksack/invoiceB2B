export default defineNuxtRouteMiddleware((to) => {
    const token = useCookie('token').value;
    const role = useCookie('role').value;
    const publicPages = ['/', '/login', '/register', '/2fa'];

    if (publicPages.includes(to.path)) return;

    if (!token) return navigateTo('/login');

    if (to.path.startsWith('/admin') && role !== 'admin') return navigateTo('/home');

    if (to.path.startsWith('/home') && role === 'admin') return navigateTo('/admin');
});


const logout = () => {
    const tokenCookie = useCookie('token');
    const roleCookie = useCookie('role');

    tokenCookie.value = null;
    roleCookie.value = null;

    navigateTo('/login');
};
