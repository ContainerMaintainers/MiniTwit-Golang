<script lang="ts">

    import jwt_decode from "jwt-decode";
    import { cookie } from '../stores/cookieStore';

    export let data: any
    $: ({ cookieValue } = data);

    cookie.subscribe(value => {
        cookieValue = value
    })

    interface jwt_token{
        exp: number
        sub: number
    }
    
    let decoded:jwt_token = jwt_decode(cookieValue!)
    console.log(decoded.exp) 
    console.log(Math.floor(Date.now()/1000))

    const logout = async () => {
        console.log("logout func")
        await fetch('http://localhost:8080/logout', {
			method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            }
		})
        window.location.href = '/'
    }

</script>

<div class=" w-screen h-32 bg-cyan-500 flex items-center">

    <h1 class=" ml-16 text-xl ">MiniTwit</h1>

    <div class=" justify-end w-full mr-10 flex gap-7">
        <button on:click={() => window.location.href = '/'}>Home</button>
        {#if decoded.exp > (Math.floor(Date.now()/1000))} <!-- divide by 1k bcs for some reason date.now() is longer -->
            
            <button on:click={() => window.location.href = '/register'}>register</button>
            <button on:click={() => window.location.href = '/login'} >login</button>
        {:else}
            <button on:click={() => logout} >logout</button>
        {/if}
        
    </div>
</div>