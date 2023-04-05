<script>
    import { load } from './+page'
    load()
</script>

<div class=" w-screen h-32 bg-cyan-500 flex items-center">

    <h1 class=" ml-16 text-xl ">MiniTwit</h1>

    <div class=" justify-end w-full mr-10 flex gap-7">
        <button class=" text-white" on:click={() => window.location.href = '/'}>Home</button>
        <button on:click={() => window.location.href = '/register'}>register</button>
        <button on:click={() => window.location.href = '/login'} >login</button>
    </div>
</div>