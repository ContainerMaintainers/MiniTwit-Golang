
import { cookie } from '../stores/cookieStore';
import type { LayoutServerLoad, LayoutServerLoadEvent } from './$types';

export const load: LayoutServerLoad = async (event: LayoutServerLoadEvent) => {

    const getToken = () => {
        let jwt_token = event.cookies.get('UserAuthorization');
        console.log('LOAD', jwt_token);
        cookie.update(_ => jwt_token!);
        return jwt_token
    }

    
    return {
        jwtToken: await getToken()
    }

}