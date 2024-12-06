# Simple wrapper for the [manifoldco/promptui](https://github.com/manifoldco/promptui) package

## Goal

I needed to ask a few questions in the console and then get the answers for some logic. promptui does this, but it seemed to me that it lacked a bit of "sugar".

This package does exactly that

### Example

Ask age

```
func main() {
	prompt := promptui.Prompt{
		Label: "How old are you?",
		Validate: func(input string) error {
			_, err := strconv.ParseInt(input, 10, 64)
			if err != nil {
				return errors.New("please enter a valid number")
			}
			return nil
		},
		Default: "18",
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	age,err := strconv.Atoi(result)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if age < 18 {
		fmt.Println("You are not allowed to enter")
	} else {
		fmt.Println("You are allowed to enter")
	}
}
```

With this package, you can do the same thing MORE EASILY, like this:

```
func main() {
	prompt := New()
	age, err := prompt.AskInt("Enter your name",18)
	if err != nil {
		fmt.Println("Prompt error: %v", err)
		return
	}

	if age < 18 {
		fmt.Println("You are not allowed to enter")
	} else {
		fmt.Println("You are allowed to enter")
	}
}
```

That's it