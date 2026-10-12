package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Notifier interface {
	Notify(msg string) error
}

type Email struct {
	Address string
}

func (e Email) Notify(msg string) error {
	fmt.Printf("Sending Email to %s: %s\n", e.Address, msg)
	return nil
}

type SMS struct {
	Phone string
}

func (s SMS) Notify(msg string) error {
	if len(msg) > 20 {
		return fmt.Errorf("SMS error: message too long (>20 chars)")
	}
	fmt.Printf("Sending SMS to %s: %s\n", s.Phone, msg)
	return nil
}

type Telegram struct {
	Chat   string
	SentCt int
}

func (t *Telegram) Notify(msg string) error {
	t.SentCt++
	fmt.Printf("Sending Telegram to %s (Msg #%d): %s\n", t.Chat, t.SentCt, msg)
	return nil
}

var _ Notifier = Email{}
var _ Notifier = SMS{}
var _ Notifier = (*Telegram)(nil)

func Broadcast(ns []Notifier, msg string) []error {
	var errs []error
	for _, n := range ns {
		if err := n.Notify(msg); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

type StudentGrade struct {
	Name  string
	Grade int
}

func writeGradeTable(w io.Writer, grades []StudentGrade) {
	fmt.Fprintln(w, "Student Name   | Grade |")
	for _, g := range grades {
		fmt.Fprintf(w, "| %-14s | %5d |\n", g.Name, g.Grade)
	}
}

func main() {
	notifiers := []Notifier{
		Email{Address: "alex@kbtu.kz"},
		SMS{Phone: "+77777777777"},
		&Telegram{Chat: "@jst_testing"},
	}

	errs := Broadcast(notifiers, "Hello team!")
	if len(errs) > 0 {
		fmt.Printf("Encountered %d errors during broadcast.\n\n", len(errs))
	}

	grades := []StudentGrade{
		{"Aruzhan", 95},
		{"Dias", 80},
		{"Timur", 88},
	}

	fmt.Println("Output to Terminal (os.Stdout)")
	writeGradeTable(os.Stdout, grades)

	fmt.Println("\nOutput to Strings.Builder")
	var sb strings.Builder
	writeGradeTable(&sb, grades)
	fmt.Print(sb.String())
}
