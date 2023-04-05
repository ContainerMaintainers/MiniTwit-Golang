<script lang="ts">

    const register = async (username: string, email: string, password: string, password2: string) => {
        console.log(username, email, password, password2)
        let hey = JSON.stringify({
				username:username,
				email:email,
                password:password,
                password2:password2
			})
        console.log(hey)
        await fetch('http://localhost:8080/register', {
			method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
			body: hey
		})
    }

    let username = ''
    let email = ''
    let pw = ''
    let pw2 = ''

</script>

<div class="">

    <div class=" mx-auto mt-28 h-80 w-72 bg-slate-600 rounded-md form-control ">
        
        <input bind:value={username} type="text" placeholder="Username" class="input input-bordered w-full max-w-xs" />
        <input bind:value={email} type="text" placeholder="Email" class="input input-bordered w-full max-w-xs" />
        <input bind:value={pw} type="password" placeholder="Password" class="input input-bordered w-full max-w-xs" />
        <input bind:value={pw2} type="password" placeholder="Password" class="input input-bordered w-full max-w-xs" />
        <button class="btn" on:click={() => register(username, email, pw, pw2)}>Submit</button>
    </div>

</div>