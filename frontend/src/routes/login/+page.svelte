<script lang="ts">

    const login = async (username: string, password: string) => {
        console.log("hej")
        await fetch('http://localhost:8080/login', {
			method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
			body: JSON.stringify({
				username:username,
				password:password
			})
		})
        window.location.href = '/'
    }

    let username = ''
    let password = ''

</script>

<div class="">

    <div class=" mx-auto mt-28 h-80 w-72 bg-slate-600 rounded-md form-control ">
        
        <input bind:value={username} type="text" placeholder="Username" class="input input-bordered w-full max-w-xs" />
        <input bind:value={password} type="password" placeholder="Password" class="input input-bordered w-full max-w-xs" />
        <button class="btn" on:click={() => login(username, password)}>Submit</button>
    </div>

</div>