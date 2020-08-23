# Image analyzer 

### Description: 
Program to analyze number of faces in images. It takes path to a directory containing images as the first argument. It depends on *Microsoft Coginitve Services*.

### Env variables: 

* "MSCV_ENDPOINT" : Microsoft Coginitive Services endpoint
* "MSCV_SUBKEY" : Corresponding Microsoft cognitive services subscription key.

Note: 5 sec delay between requests is added, because of the limits imposed by Microsoft services.