#!/bin/bash

# Скрипт для автоматического обновления покрытия тестами в README.md
# Использование: ./scripts/update_coverage.sh

set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧪 Запуск тестов и генерация отчета о покрытии...${NC}"

# Создаем директорию для отчетов если её нет
mkdir -p coverage

# Запускаем тесты с генерацией покрытия
go test -coverprofile=coverage/coverage.out ./... 2>/dev/null || {
    echo -e "${RED}❌ Ошибка при запуске тестов${NC}"
    exit 1
}

# Получаем общее покрытие
COVERAGE=$(go tool cover -func=coverage/coverage.out | tail -1 | awk '{print $3}' | sed 's/%//')

if [ -z "$COVERAGE" ]; then
    echo -e "${RED}❌ Не удалось получить данные о покрытии${NC}"
    exit 1
fi

echo -e "${GREEN}📊 Общее покрытие тестами: ${COVERAGE}%${NC}"

# Генерируем HTML отчет
go tool cover -html=coverage/coverage.out -o coverage/coverage.html
echo -e "${BLUE}📄 HTML отчет сохранен в coverage/coverage.html${NC}"

# Определяем цвет бейджа в зависимости от покрытия
COVERAGE_NUM=$(echo $COVERAGE | cut -d. -f1)
if [ "$COVERAGE_NUM" -ge 80 ]; then
    BADGE_COLOR="brightgreen"
elif [ "$COVERAGE_NUM" -ge 60 ]; then
    BADGE_COLOR="yellow"
elif [ "$COVERAGE_NUM" -ge 40 ]; then
    BADGE_COLOR="orange"
else
    BADGE_COLOR="red"
fi

# Создаем бейдж для покрытия
COVERAGE_BADGE="![Coverage](https://img.shields.io/badge/coverage-${COVERAGE}%25-${BADGE_COLOR})"

# Обновляем README.md
README_FILE="README.md"

if [ ! -f "$README_FILE" ]; then
    echo -e "${RED}❌ Файл README.md не найден${NC}"
    exit 1
fi

# Создаем временный файл
TEMP_FILE=$(mktemp)

# Флаг для отслеживания, был ли найден раздел с покрытием
COVERAGE_SECTION_FOUND=false

while IFS= read -r line; do
    # Проверяем, начинается ли строка с бейджа покрытия
    if [[ $line =~ ^\!\[Coverage\] ]]; then
        echo "$COVERAGE_BADGE" >> "$TEMP_FILE"
        COVERAGE_SECTION_FOUND=true
        echo -e "${GREEN}✅ Обновлен бейдж покрытия в README.md${NC}"
    # Проверяем раздел со статистикой покрытия
    elif [[ $line =~ ^##.*[Пп]окрытие.*[Тт]естами ]]; then
        echo "$line" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
        echo "**Общее покрытие:** ${COVERAGE}%" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
        echo "📊 **Детализированное покрытие по модулям:**" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
        echo "\`\`\`" >> "$TEMP_FILE"
        go tool cover -func=coverage/coverage.out | grep -v "total:" | while read -r module_line; do
            echo "$module_line" >> "$TEMP_FILE"
        done
        echo "\`\`\`" >> "$TEMP_FILE"
        echo "" >> "$TEMP_FILE"
        echo "*Отчет автоматически обновлен $(date '+%Y-%m-%d %H:%M:%S')*" >> "$TEMP_FILE"

        # Пропускаем старое содержимое раздела до следующего заголовка
        while IFS= read -r next_line; do
            if [[ $next_line =~ ^##.* ]] || [[ $next_line =~ ^#.* ]]; then
                echo "$next_line" >> "$TEMP_FILE"
                break
            fi
        done
        COVERAGE_SECTION_FOUND=true
        echo -e "${GREEN}✅ Обновлен раздел с покрытием в README.md${NC}"
    else
        echo "$line" >> "$TEMP_FILE"
    fi
done < "$README_FILE"

# Если раздел покрытия не найден, добавляем его в конец
if [ "$COVERAGE_SECTION_FOUND" = false ]; then
    echo "" >> "$TEMP_FILE"
    echo "## 🧪 Покрытие тестами" >> "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
    echo "$COVERAGE_BADGE" >> "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
    echo "**Общее покрытие:** ${COVERAGE}%" >> "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
    echo "📊 **Детализированное покрытие по модулям:**" >> "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
    echo "\`\`\`" >> "$TEMP_FILE"
    go tool cover -func=coverage/coverage.out | grep -v "total:" >> "$TEMP_FILE"
    echo "\`\`\`" >> "$TEMP_FILE"
    echo "" >> "$TEMP_FILE"
    echo "*Отчет автоматически обновлен $(date '+%Y-%m-%d %H:%M:%S')*" >> "$TEMP_FILE"
    echo -e "${GREEN}✅ Добавлен новый раздел с покрытием в README.md${NC}"
fi

# Заменяем оригинальный файл
mv "$TEMP_FILE" "$README_FILE"

echo -e "${GREEN}🎉 Покрытие тестами успешно обновлено!${NC}"
echo -e "${YELLOW}💡 Для просмотра детального отчета откройте: coverage/coverage.html${NC}"

# Выводим краткую статистику
echo -e "${BLUE}📈 Краткая статистика:${NC}"
echo "- Общее покрытие: ${COVERAGE}%"
echo "- HTML отчет: coverage/coverage.html"
echo "- Обновлен файл: README.md"
