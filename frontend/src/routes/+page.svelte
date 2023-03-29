<script>
    
    //let data = fetch('localhost:8080/').then(res => res.json())

    let apiResJson = fetch('http://localhost:8080/public').then(res => res.json())
    let data = apiResJson.then(res => res['messages'])
</script>

<h1 class=" bg-green-700 font-bold" >tweets:</h1>
{#await data}
    <div class="bg-green-200">
        <h1 class="font-bold">Loading...</h1>
    </div>
{:then data}
{#each  data as tweet}
    <div class="bg-green-200">
        <h1 class="font-bold">{tweet.Text}</h1>
    </div>
{/each}
{/await}