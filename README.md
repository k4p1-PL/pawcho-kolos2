# PACO - Zadanie 2

## Opis Łańcucha GHAction
Łańcuch (pipeline) został zaprojektowany w celu zautomatyzowania budowy, skanowania i publikacji obrazu aplikacji. Działa on według poniższego schematu:
1. Konfiguruje środowisko (QEMU, Buildx) umożliwiające budowanie wieloarchitektoniczne (`linux/amd64`, `linux/arm64`).
2. Buduje obraz testowy i ładuje go lokalnie z wykorzystaniem dostępnego cache'a.
3. **Test CVE:** Uruchamia skaner **Trivy**. Został on wybrany jako narzędzie optymalne dla środowisk CI/CD, ponieważ natywnie wspiera blokowanie procesu (`exit-code: '1'`) w przypadku wykrycia podatności `CRITICAL` lub `HIGH`.
4. Jeśli test zostanie zaliczony pomyślnie, obraz jest budowany na obydwie architektury i wysyłany do GitHub Container Registry (ghcr.io).

## System tagowania i przechowywania Cache (Z Uzasadnieniem)

### 1. Tagowanie danych Cache
* **Format:** `[dockerhub-user]/paco-cache:buildcache`
* **Uzasadnienie:** Zgodnie z dobrymi praktykami oraz wymaganiami zadania, dane z cache'a procesu budowania (`type=registry`, `mode=max`) są wypychane do dedykowanego, oddzielnego repozytorium na DockerHub (`paco-cache`). Separacja warstw cache od repozytorium docelowego z obrazem to zalecana praktyka (tzw. zjawisko *Cache Pollution*). Pozwala to utrzymać porządek w tagach aplikacji produkcyjnej – użytkownicy końcowi widzą tylko gotowe obrazy, a nie śmieciowe warstwy buildera.

### 2. Tagowanie obrazu docelowego (ghcr.io)
* **Format:** Obraz tagowany jest podwójnie przy każdym poprawnym procesie CI:
  * `latest` - znacznik ruchomy, nadpisywany przy każdym commicie.
  * `${{ github.sha }}` - znacznik niezmienny (immutable), oparty o sumę kontrolną commita w Git.
* **Uzasadnienie:** Użycie podwójnego tagowania to standard w inżynierii DevOps. Tag `latest` służy wygodzie (np. dla szybkiego testowania uruchomieniowego najnowszej wersji). Natomiast tag oparty na SHA commita gwarantuje **identyfikowalność (traceability)** i niezmienność. Umozliwia to precyzyjne określenie, z jakiej wersji kodu źródłowego powstał dany obraz kontenera, co jest kluczowe w debugowaniu i wycofywaniu zmian na środowiskach produkcyjnych.